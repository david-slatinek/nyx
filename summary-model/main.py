from transformers import pipeline

if __name__ == "__main__":
    summarizer = pipeline("summarization", model="philschmid/bart-large-cnn-samsum")

    conversation = '''Jeff: Can I train a 🤗 Transformers model on Amazon SageMaker? 
    Philipp: Sure you can use the new Hugging Face Deep Learning Container. 
    Jeff: ok.
    Jeff: and how can I get started? 
    Jeff: where can I find documentation? 
    Philipp: ok, ok you can find everything here.'''

    summary = summarizer(conversation)

    print(summary[0]["summary_text"])
