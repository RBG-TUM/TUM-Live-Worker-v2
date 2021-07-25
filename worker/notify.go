package worker

import (
	"context"
	"fmt"
	"github.com/joschahenningsen/TUM-Live-Worker-v2/cfg"
	"github.com/joschahenningsen/TUM-Live-Worker-v2/pb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func notifyStreamStart(streamCtx *StreamContext) {
	client, conn, err := GetClient()
	if err != nil {
		log.WithError(err).Error("Unable to dial server")
		return
	}
	resp, err := client.NotifyStreamStarted(context.Background(), &pb.StreamStarted{
		WorkerID:   cfg.WorkerID,
		StreamID:   streamCtx.streamId,
		HlsUrl:     fmt.Sprintf("https://live.stream.lrz.de/livetum/%s{{quality}}/playlist.m3u8?dvr", streamCtx.getStreamName()),
		SourceType: streamCtx.streamVersion,
	})
	if err != nil || !resp.Ok {
		log.WithError(err).Error("Could not notify stream finished")
	}
	_ = conn.Close()
}

func notifyStreamDone(streamID uint32) {
	client, conn, err := GetClient()
	if err != nil {
		log.WithError(err).Error("Unable to dial server")
		return
	}
	resp, err := client.NotifyStreamFinished(context.Background(), &pb.StreamFinished{
		WorkerID: cfg.WorkerID,
		StreamID: streamID,
	})
	if err != nil || !resp.Ok {
		log.WithError(err).Error("Could not notify stream finished")
	}
	_ = conn.Close()
}

func notifyTranscodingDone(streamCtx *StreamContext) {
	client, conn, err := GetClient()
	if err != nil {
		log.WithError(err).Error("Unable to dial server")
		return
	}
	resp, err := client.NotifyTranscodingFinished(context.Background(), &pb.TranscodingFinished{
		WorkerID: cfg.WorkerID,
		StreamID: streamCtx.streamId,
		FilePath: streamCtx.getTranscodingFileName(),
	})
	if err != nil || !resp.Ok {
		log.WithError(err).Error("Could not notify stream finished")
	}
	_ = conn.Close()
}

func notifyUploadDone(streamCtx *StreamContext) {
	client, conn, err := GetClient()
	if err != nil {
		log.WithError(err).Error("Unable to dial server")
		return
	}
	resp, err := client.NotifyUploadFinished(context.Background(), &pb.UploadFinished{
		WorkerID:   cfg.WorkerID,
		StreamID:   streamCtx.streamId,
		HLSUrl:     fmt.Sprintf("https://stream.lrz.de/vod/_definst_/mp4:tum/RBG/%s.mp4/playlist.m3u8", streamCtx.getStreamName()),
		SourceType: "",
	})
	if err != nil || !resp.Ok {
		log.WithError(err).Error("Could not notify upload finished")
	}
	_ = conn.Close()
}

func GetClient() (pb.FromWorkerClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:50052", cfg.MainBase), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return pb.NewFromWorkerClient(conn), conn, nil
}