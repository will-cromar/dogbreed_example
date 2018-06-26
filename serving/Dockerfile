FROM debian

RUN apt-get update && apt-get -y install curl gnupg

# From TF Serving instructions
RUN echo "deb [arch=amd64] http://storage.googleapis.com/tensorflow-serving-apt stable tensorflow-model-server tensorflow-model-server-universal" | tee /etc/apt/sources.list.d/tensorflow-serving.list
RUN curl https://storage.googleapis.com/tensorflow-serving-apt/tensorflow-serving.release.pub.gpg | apt-key add -
RUN apt-get update && apt-get install tensorflow-model-server

COPY model /model

CMD tensorflow_model_server --model_base_path=/model --rest_api_port=$PORT