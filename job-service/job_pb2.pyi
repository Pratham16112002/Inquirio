from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class ParseResumeRequest(_message.Message):
    __slots__ = ("resume_file_content", "file_name")
    RESUME_FILE_CONTENT_FIELD_NUMBER: _ClassVar[int]
    FILE_NAME_FIELD_NUMBER: _ClassVar[int]
    resume_file_content: bytes
    file_name: str
    def __init__(self, resume_file_content: _Optional[bytes] = ..., file_name: _Optional[str] = ...) -> None: ...

class ParseResumeResponse(_message.Message):
    __slots__ = ("job_titles", "skills", "experience")
    JOB_TITLES_FIELD_NUMBER: _ClassVar[int]
    SKILLS_FIELD_NUMBER: _ClassVar[int]
    EXPERIENCE_FIELD_NUMBER: _ClassVar[int]
    job_titles: _containers.RepeatedScalarFieldContainer[str]
    skills: _containers.RepeatedScalarFieldContainer[str]
    experience: int
    def __init__(self, job_titles: _Optional[_Iterable[str]] = ..., skills: _Optional[_Iterable[str]] = ..., experience: _Optional[int] = ...) -> None: ...

class CalculateRelevancyRequest(_message.Message):
    __slots__ = ("resume_skills", "resume_experience", "job_description")
    RESUME_SKILLS_FIELD_NUMBER: _ClassVar[int]
    RESUME_EXPERIENCE_FIELD_NUMBER: _ClassVar[int]
    JOB_DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    resume_skills: _containers.RepeatedScalarFieldContainer[str]
    resume_experience: str
    job_description: str
    def __init__(self, resume_skills: _Optional[_Iterable[str]] = ..., resume_experience: _Optional[str] = ..., job_description: _Optional[str] = ...) -> None: ...

class CalculateRelevancyResponse(_message.Message):
    __slots__ = ("relevancy_score",)
    RELEVANCY_SCORE_FIELD_NUMBER: _ClassVar[int]
    relevancy_score: float
    def __init__(self, relevancy_score: _Optional[float] = ...) -> None: ...
