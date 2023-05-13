from transformers import pipeline

if __name__ == "__main__":
    classifier = pipeline("zero-shot-classification", model="facebook/bart-large-mnli")

    sequence_to_classify = "I really want to buy a new laptop, but I am not sure which one to buy."
    candidate_labels = ["pc", "laptops", "Dell", "Lenovo", "Apple"]

    res = classifier(sequence_to_classify, candidate_labels, multi_label=True)

    print(res)
