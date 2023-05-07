import schema_pb2_grpc as pb2_grpc
import schema_pb2 as pb2
from concurrent import futures
import grpc


class DialogServiceServicer(pb2_grpc.DialogServiceServicer):
    def Dialog(self, request, context):
        message = request.text
        result = f'Hello I am up and running received "{message}" message from you'

        return pb2.DialogResponse(answer=result)


if __name__ == "__main__":
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_DialogServiceServicer_to_server(DialogServiceServicer(), server)
    server.add_insecure_port("[::]:9080")
    server.start()
    server.wait_for_termination()
