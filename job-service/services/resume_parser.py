import os
import json
import logging
from typing import List
import grpc
import requests
import re
from pydantic import BaseModel, Field
from utils.text_cleaner import clean_resume_text

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# --- Pydantic Model for Structured Output ---
class ResumeData(BaseModel):
    """Structured data extracted from a resume."""
    job_titles: List[str] = Field(default_factory=list, description="Extracted job titles.")
    skills: List[str] = Field(default_factory=list, description="Extracted skills.")
    experience: int = Field(default=0, description="A summary of the work experience.")

# --- Resume Parser using Local Ollama ---
class ResumeParser:
    def __init__(self, model_name: str = "llama3.1:latest"):
        """
        Uses a local Ollama model to extract structured data from resumes.
        """
        self.model_name = model_name
        self.api_url = "http://localhost:11434/api/generate"
        logger.info(f"Ollama model set to: {self.model_name}")

    def _query_ollama(self, prompt: str) -> str:
        """
        Sends a request to the local Ollama server.
        """
        try:
            payload = {
                "model": self.model_name,
                "prompt": prompt,
                "stream": False
            }

            response = requests.post(self.api_url, json=payload)
            response.raise_for_status()

            # Ollama returns: {"model":"...","created_at":"...","response":"..."}
            data = response.json()
            return data.get("response", "")

        except Exception as e:
            logger.error(f"Error communicating with Ollama: {e}", exc_info=True)
            return ""

    def parse(self, resume_text: str, context=None) -> ResumeData:
        if not resume_text:
            logger.warning("Input resume text is empty.")
            if context:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("File has nothing to process")
            return None

        cleaned_context = clean_resume_text(resume_text)
        logger.info("Extracting information from resume using Ollama...")

        prompt = f"""
        You are a resume parser AI. 
        Extract the following information from the resume below and return valid JSON only (no extra text):

        Resume:
        {resume_text}

        Return the result strictly in JSON format with these keys:
        {{
            "job_titles": [list of job titles that you can extract from this resume ( NOT object or array and keep the array empty if no job title is found)],
            "skills": [list of skills that you can extract from this resume every where in the resume ( NOT object or array and keep the array empty if no skill is found)],
            "experience": "Total experience in years carefully extract the work experience section and add one the duration of each of them ( should be 32 bit integer only add 0 if experience is in months take round of the number)"
        }}
        """
        llm_output = self._query_ollama(prompt)

        raw = llm_output.strip()
        match = re.search(r"\{[\s\S]*\}", raw)
        if not match:
            logger.error(f"Ollama output not valid JSON:\n{llm_output}")
            if context:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details("No JSON object found in LLM output.")
                return None
        json_str = match.group(0)
        try:
            parsed_json = json.loads(json_str)
            job_titles = parsed_json.get("job_titles", [])
            skills = parsed_json.get("skills", [])
            experience = parsed_json.get("experience", "")
        except json.JSONDecodeError:
            logger.error(f"Ollama output not valid JSON:\n{llm_output}")
            if context:
                context.set_code(grpc.StatusCode.INTERNAL)
                context.set_details("LLM output not valid JSON.")
            return None

        extracted_data = {
            "job_titles": job_titles,
            "skills": skills,
            "experience": experience
        }
        return ResumeData(
            job_titles=extracted_data["job_titles"],
            skills=extracted_data["skills"],
            experience=extracted_data["experience"],
        )


# --- Example Usage ---
if __name__ == '__main__':
    sample_resume = """
    John Doe
    Software Engineer

    Experience:
    - Senior Developer at Tech Corp (2020-Present)
      - Led a team to build microservices using Python, Django, and Kubernetes.
    - Software Engineer at Innovate LLC (2018-2020)
      - Built web applications with React and Node.js.

    Skills:
    Python, Java, JavaScript, React, Node.js, Docker, Kubernetes, SQL, Git.
    """

    print("Initializing ResumeParser (Ollama)...")
    parser = ResumeParser(model_name="llama3.1:latest")

    print("\n--- Parsing Sample Resume ---")
    parsed_data = parser.parse(sample_resume)

    print("\n--- Extracted Data ---")
    if parsed_data:
        print(f"Job Titles: {parsed_data.job_titles}")
        print(f"Skills: {parsed_data.skills}")
        print(f"Experience Summary: {parsed_data.experience}")
