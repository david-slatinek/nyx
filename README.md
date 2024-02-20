<div align="center">
    <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>
    <img alt="AWS" src="https://img.shields.io/badge/Amazon_AWS-232F3E?style=for-the-badge&logo=amazonaws&logoColor=white"/>
    <img alt="AWS Lambda" src="https://img.shields.io/badge/aws_lambda-FF9900?style=for-the-badge&logo=awslambda&logoColor=white"/>
    <img alt="MySQL" src="https://img.shields.io/badge/MySQL-005C84?style=for-the-badge&logo=mysql&logoColor=white"/>
    <img alt="Docker" src="https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white"/>
    <img alt="Amazon EC2" src="https://img.shields.io/badge/amazon_ec2-FF9900?style=for-the-badge&logo=amazonec2&logoColor=white"/>
    <img alt="Github Actions" src="https://img.shields.io/badge/Github%20Actions-282a2e?style=for-the-badge&logo=githubactions&logoColor=367cfe"/>
    <img alt="Amazon SQS" src="https://img.shields.io/badge/amazon_sqs-FF4F8B?style=for-the-badge&logo=amazonsqs&logoColor=white"/>
    <img alt="React" src="https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB"/>
    <img alt="dvc" src="https://img.shields.io/badge/dvc-13ADC7?style=for-the-badge&logo=dvc&logoColor=white"/>
    <img alt="Python" src="https://img.shields.io/badge/python-3776AB?style=for-the-badge&logo=python&logoColor=white"/>
    <img alt="Apache CouchDB" src="https://img.shields.io/badge/apache_couchdb-E42528?style=for-the-badge&logo=apachecouchdb&logoColor=white"/>
    <img alt="Grafana" src="https://img.shields.io/badge/grafana-F46800?style=for-the-badge&logo=grafana&logoColor=white"/>
</div>

![Logo](images/logo.png)

# Project description

A recommendation system for recommending products for an online store that consists of the following components:

- [Web application](website)
    - React application that allows users to chat with a chatbot.
- [Dialog API](dialog-api)
    - API for storing dialogs in the database and responding to the user.
- [Dialog model](dialog-model)
    - Machine learning model that generates responses based on the user's input.
- [CouchDB](couchdb)
    - Database for storing dialogs, summaries, and product categories.
- [Summary API](summary-api)
    - API for storing chat summaries in the database.
- [Summary model](summary-model)
    - Machine learning model that generates chat summary.
- [Categories API](categories-api)
    - API for managing product categories.
- [Recommend API](recommend-api)
    - API for recommending products to the user.
- [Recommend model](recommend-model)
    - Machine learning model that recommends products for the user based on the chat.
- [Recommend db](recommend-db)
    - Database for storing recommendations.
- [Recommend email](recommend-lambda)
    - Lambda function that sends a recommendation email to the user.
- [Tracking API](tracking-api)
    - API for tracking user clicks in the email.
- [GitHub Actions](.github/workflows)
    - CI/CD pipelines for the project.
- [Grafana dashboard](monitoring)
    - Dashboard for project monitoring.

AWS services were mainly used for the project infrastructure:

- Amazon SQS
    - Message queue for communication between the components.
- Amazon EC2
    - Virtual machine for hosting the project.
- Amazon ECR
    - Docker container registry for storing Docker images.
        - Due to the free tier limitations, Docker Hub was used for storing images for models.
- Amazon Lambda
    - Serverless computing service for running the recommendation email function.

Due to the complexity of the project, DVC and DagsHub were used for managing the machine learning models and data.

# Project infrastructure

<div align="center">
  <img src="images/diagram.png" alt="Project infrastructure">
  <br/>
  <i>Project infrastructure.</i>
</div>
<br/>

# Website

<div align="center">
  <img src="images/website.png" alt="Website">
  <br/>
  <i>Website.</i>
</div>

# Recommendation email

<div align="center">
  <img src="images/recommendation.png" alt="Recommendation email">
  <br/>
  <i>Recommendation email.</i>
</div>

# Monitoring

<div align="center">
  <img src="images/monitoring.png" alt="Monitoring">
  <br/>
  <i>Monitoring.</i>
</div>
