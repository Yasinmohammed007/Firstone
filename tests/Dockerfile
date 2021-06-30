FROM python:3.9

RUN mkdir /code
COPY ./ /code/
# COPY recommendations/ /service/recommendations/
WORKDIR /code/test
RUN python3 -m pip install --upgrade pip
RUN python3 -m pip install -r ../requirements.txt
# RUN python3 -m grpc_tools.protoc -I ../protobufs --python_out=. --grpc_python_out=. ../protobufs/recommendations.proto
ENTRYPOINT [ "pytest", "sample-test.py" ]