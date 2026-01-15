from ordered_set import OrderedSet
import asyncio

class NoqueueMessage:
    def __init__(self, action, data):
        self.action_requirements = {
            "PUB": OrderedSet(["topic", "message"]),
            "SUB": OrderedSet(["topic"]),
            "EXIT": OrderedSet([])
        }
        if action not in self.action_requirements:
            raise Exception(f"Invalid action: {action}")
        if data.keys() != self.action_requirements[action]:
            raise Exception(f"Insufficient or improper keys for given action, required: {self.action_requirements[action]}")
        self.action, self.data = action, data

    def format(self):
        string = f"{self.action} {' '.join(self.data[key] for key in self.action_requirements[self.action])}\n"
        return string.encode()

class NoqueueClient:
    def __init__(self, address):
        self.host, self.port = address
        self.connected = False
        
    async def create_connection(self):
        if self.connected:
            return
        self.reader, self.writer = await asyncio.open_connection(self.host, self.port)
        print("CONNECTED")
        self.connected = True

    async def send(self, message):
        message = message.format()
        if not self.connected:  
            raise Exception(f"Sending message {message} on a disconnected client")
        self.writer.write(message)
        await self.writer.drain()

    async def exit(self):
        exit_message = NoqueueMessage("EXIT", {})
        await self.send(exit_message)
        self.writer.close()
        await self.writer.wait_closed()
        self.connected = False


class NoqueueSubscriber(NoqueueClient):
    def __init__(self, address):
        super().__init__(address)

    async def subscribe(self, topic):
        await self.create_connection()
        sub_message = NoqueueMessage("SUB", {"topic": topic})
        await self.send(sub_message)

    async def wait_for_message(self):
        if not self.connected:
            raise Exception("Waiting for message on an unconnected NoqueueSubscriber")
        while True:
            data = await self.reader.readline()
            if not data:
                break
            yield data.decode().rstrip("\n")
        
    def __aiter__(self):
        return self.wait_for_message()


class NoqueuePublisher(NoqueueClient):
    def __init__(self, address):
        super().__init__(address)

    async def publish(self, topic, message):
        await self.create_connection()
        pub_message = NoqueueMessage("PUB", {"topic": topic, "message": message})
        await self.send(pub_message)
        await self.exit()

    