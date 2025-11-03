from concurrent import futures
import grpc
import job_pb2
import job_pb2_grpc
import logging
from services.resume_parser import ResumeParser

logging.basicConfig(level=logging.INFO)

class JobServiceServicer(job_pb2_grpc.JobServiceServicer):
    def __init__(self):
        logging.info("Initializing JobServiceServicer...")
        self.parser = ResumeParser()

    def GetJobs(self, request, context):

        # Parsing the resume
        parsed_data = self.parser.parse(request.resume_text)

        return job_pb2.ParseResumeResponse(
            job_titles=parsed_data.job_titles,
            skills=parsed_data.skills,  
            experience=parsed_data.experience
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    job_pb2_grpc.add_JobServiceServicer_to_server(AIServiceServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started, listening on port 50051.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
