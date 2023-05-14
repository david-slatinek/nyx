from concurrent import futures

import grpc
from transformers import pipeline

import schema_pb2 as pb2
import schema_pb2_grpc as pb2_grpc


class RecommendServiceServicer(pb2_grpc.RecommendServiceServicer):
    def RecommendDialog(self, request, context):
        categories = request.categories

        labels = []
        scores = []
        for dialog in request.dialogs:
            res = classifier(dialog, categories, multi_label=True)
            labels.append(res["labels"])
            scores.append(res["scores"])

        res = []
        for i in range(len(labels)):
            res.append(pb2.RecommendResponse(
                text=request.dialogs[i],
                labels=labels[i],
                scores=scores[i]
            ))

        return pb2.RecommendResponseList(
            responses=res
        )

    def RecommendSummary(self, request, context):
        categories = request.categories

        res = classifier(request.summary, categories, multi_label=True)

        return pb2.RecommendResponse(
            text=request.summary,
            labels=res["labels"],
            scores=res["scores"]
        )


if __name__ == "__main__":
    classifier = pipeline("zero-shot-classification", model="facebook/bart-large-mnli")

    print("Model loaded")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    pb2_grpc.add_RecommendServiceServicer_to_server(RecommendServiceServicer(), server)
    server.add_insecure_port("[::]:9040")
    server.start()
    server.wait_for_termination()
