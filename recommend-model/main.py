from concurrent import futures

import grpc
from transformers import pipeline

import schema_pb2 as pb2
import schema_pb2_grpc as pb2_grpc


class RecommendServiceServicer(pb2_grpc.RecommendServiceServicer):
    def Recommend(self, request, context):
        categories = request.categories

        labels = []
        scores = []
        dialogs = []
        for dialog in request.dialogs:
            res = classifier(dialog, categories, multi_label=True)
            dialogs.append(dialog)
            labels.append(res["labels"])
            scores.append(res["scores"])

        res = classifier(request.summary, categories, multi_label=True)
        dialogs.append(request.summary)
        labels.append(res["labels"])
        scores.append(res["scores"])

        res = []
        for i in range(len(labels)):
            res.append(pb2.RecommendResponse(
                dialog=dialogs[i],
                labels=labels[i],
                scores=scores[i]
            ))

        return pb2.RecommendResponseList(
            responses=res
        )


if __name__ == "__main__":
    classifier = pipeline("zero-shot-classification", model="facebook/bart-large-mnli")

    print("Model loaded")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    pb2_grpc.add_RecommendServiceServicer_to_server(RecommendServiceServicer(), server)
    server.add_insecure_port("[::]:9040")
    server.start()
    server.wait_for_termination()

    # sequence_to_classify = "I really want to buy a new laptop, but I am not sure which one to buy."
    # candidate_labels = ["pc", "laptops", "Dell", "Lenovo", "Apple"]
    # res = classifier(sequence_to_classify, candidate_labels, multi_label=True)
    # print(res)
