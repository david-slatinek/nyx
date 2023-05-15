FROM python:3.10

WORKDIR /recommend-model

COPY pyproject.toml .

RUN python -m pip install --upgrade pip
RUN pip install poetry
RUN poetry install --only main

COPY . .

CMD ["poetry", "run", "python", "main.py"]
