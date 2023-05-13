import os
from random import shuffle
import requests
from dotenv import load_dotenv

load_dotenv()

if __name__ == "__main__":
    r = requests.get(os.getenv("DIALOG_URL") + "/dialogs")
    data = r.json()

    dialogs = [dialog["text"] for dialog in data]
    shuffle(dialogs)

    test_size = int(len(dialogs) * 0.3)
    train = dialogs[:len(dialogs) - test_size]
    test = dialogs[len(dialogs) - test_size:]

    with open("../data/train.txt", "w") as f:
        for dialog in train:
            f.write(dialog + "\n")

    with open("../data/test.txt", "w") as f:
        for dialog in test:
            f.write(dialog + "\n")
