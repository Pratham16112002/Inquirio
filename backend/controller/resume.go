package controller

import (
	"Inquiro/config"
	"Inquiro/protos"
	"Inquiro/services"
	"Inquiro/utils/response"
	"io"
	"net/http"

	"google.golang.org/grpc/status"
)

type Resume struct {
	srv services.Service
	cfg config.Application
}

func (u Resume) ProcessResume(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		u.cfg.Logger.Warnw("Bad request", "error : ", err.Error())
		response.Error(w, r, "Bad request", "File too large", 400, http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("resume")
	if err != nil {
		u.cfg.Logger.Warnw("Could not find resume in request", "error : ", err.Error())
		response.Error(w, r, "Bad request", "File not found", 400, http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		u.cfg.Logger.Warnw("Could not read file", "error : ", err.Error())
		response.Error(w, r, "File unreadable", "Could not read the file content", 400, http.StatusBadRequest)
		return
	}
	u.cfg.Logger.Infow("Reading file successfull, sending to python service", "filename", header.Filename, "size", len(fileBytes))
	ctx := r.Context()
	res, err := u.cfg.Grpc.ParseResume(ctx, &protos.ParseResumeRequest{
		ResumeFileContent: fileBytes,
		FileName:          header.Filename,
	})
	if err != nil {
		u.cfg.Logger.Warnw("Could not send request to python service", "error : ", err.Error())
		st, _ := status.FromError(err)
		response.Error(w, r, "File not processed", st.Message(), int(status.Code(err)), http.StatusInternalServerError)
		return
	}
	response.Success(w, r, "File proccessed", struct {
		JobTitles  []string `json:"job_titles"`
		Skills     []string `json:"skills"`
		Experience int32    `json:"experience"`
	}{
		JobTitles:  res.JobTitles,
		Skills:     res.Skills,
		Experience: res.Experience,
	}, http.StatusOK)

}
