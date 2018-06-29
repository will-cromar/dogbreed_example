# Dog Breed Identification with TensorFlow

The purpose of this repository is to show a complete example of a project using TensorFlow, from data preparation to serving. You can see the finished product here:

https://mysterious-river-44059.herokuapp.com/ (the first load may be quite slow)

Broadly, the repository has 4 parts:
1. Data preparation
2. Model training and export
3. TensorFlow Serving container
4. Web frontend

I've included an example exported SavedModel in the serving directory, so you can skip the first two parts if you're not interested. The model is fairly lightweight and I've included Dockerfiles for both the frontend and the TF Serving backend, so you can run these just about anywhere with minimal modifications. I chose to run them on Heroku's free tier.

## Data

This repository is meant to be used with the data from Kaggle's old [Dog Breed Identification](https://www.kaggle.com/c/dog-breed-identification/) competition. It's assumed that it will be downloaded and unzipped into the `data/` directory.

## Training the Model

First, download and unzip the data as described above. Then, run `data_prep.py`, which will preprocess and turn all of the images into TFRecords. Be aware that the TFRecord file will be quite large (about 6GB). Then, run all of the code in `complete_model.ipynb`, which will train and export a transfer learning model based on MobileNet V2.

## Serving

The `serving/` directory has a Dockerfile that you can use to package the exported SavedModel. In fact, I'm fairly sure you can substitute any model for the one in the `serving/model/` directory and it will work just fine. Note that the container has to be started with the `$PORT` environment variable, which is where it will listen for REST requests. I don't use the gRPC interface at all for this project.

## Web Frontend

I've also included a very basic web frontend written in Go that provides an easy interface to send photos to the model.

This repository is still something of a work in progress. Going forward, I want to document each piece more thoroughly and perhaps do a write up of how I did this.
