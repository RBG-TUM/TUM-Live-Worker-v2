package worker

import (
	"github.com/joschahenningsen/TUM-Live-Worker-v2/cfg"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func upload(streamCtx *StreamContext) {
	log.WithField("stream", streamCtx.getStreamName()).Info("Uploading stream")
	err := post(streamCtx.getTranscodingFileName())
	if err != nil {
		log.WithField("stream", streamCtx.getStreamName()).WithError(err).Error("Error uploading stream")
	}
	log.WithField("stream", streamCtx.getStreamName()).Info("Uploaded stream")
}

func post(file string) error {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	r, w := io.Pipe()
	writer := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer writer.Close()
		err := writeFile(writer, "filename", file)
		if err != nil {
			log.Error("Cannot create form file: ", err)
			return
		}

		fields := map[string]string{
			"benutzer":    cfg.LrzUser,
			"mailadresse": cfg.LrzMail,
			"telefon":     cfg.LrzPhone,
			"unidir":      "tum",
			"subdir":      cfg.LrzSubDir,
			"info":        "",
		}

		for name, value := range fields {
			err = writeField(writer, name, value)
			if err != nil {
				log.Error("Cannot create form field: ", err)
				return
			}
		}
	}()
	rsp, err := client.Post(cfg.LrzUploadUrl, writer.FormDataContentType(), r)
	if err == nil && rsp.StatusCode != http.StatusOK {
		log.Error("Request failed with response code: ", rsp.StatusCode)
	}

	return err
}

func writeField(writer *multipart.Writer, name string, value string) error {
	formFieldWriter, err := writer.CreateFormField(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(formFieldWriter, strings.NewReader(value))
	return err
}

func writeFile(writer *multipart.Writer, fieldname string, file string) error {
	formFileWriter, err := writer.CreateFormFile(fieldname, file)
	if err != nil {
		return err
	}
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	_, err = io.Copy(formFileWriter, fileReader)
	return err
}
