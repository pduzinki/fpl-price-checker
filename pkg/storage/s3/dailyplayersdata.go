package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pduzinki/fpl-price-checker/pkg/config"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"
)

type Player struct {
	ID         int
	Name       string
	Price      int
	SelectedBy string
}

type DailyPlayersData map[int]Player

type DailyPlayersDataRepository struct {
	Uploader   *s3manager.Uploader
	Downloader *s3manager.Downloader
	Bucket     string
	Prefix     string
}

func NewDailyPlayersDataRepository(awsConfig config.AWSConfig, prefix string) (*DailyPlayersDataRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           &awsConfig.Region,
		Credentials:      credentials.NewStaticCredentials(awsConfig.ID, awsConfig.Secret, awsConfig.Token),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         &awsConfig.Endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("s3.NewDailyPlayersDataRepository failed: %w", err)
	}

	return &DailyPlayersDataRepository{
		Uploader:   s3manager.NewUploader(sess),
		Downloader: s3manager.NewDownloader(sess),
		Bucket:     awsConfig.Bucket,
		Prefix:     prefix,
	}, nil
}

func (dr *DailyPlayersDataRepository) Add(ctx context.Context, date string, players domain.DailyPlayersData) error {
	s3Players := toS3DailyPlayersData(players)

	jsonPlayers, err := json.Marshal(s3Players)
	if err != nil {
		return fmt.Errorf("s3.NewDailyPlayersDataRepository failed to mashal data: %w", err)
	}

	_, err = dr.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(dr.Bucket),
		Key:    aws.String(filepath.Join(dr.Prefix, date)),
		Body:   bytes.NewReader(jsonPlayers),
	})
	if err != nil {
		return fmt.Errorf("s3.NewDailyPlayersDataRepository failed to upload to s3: %w", err)
	}

	return nil
}

func (dr *DailyPlayersDataRepository) GetByDate(ctx context.Context, date string) (domain.DailyPlayersData, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := dr.Downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: &dr.Bucket,
		Key:    aws.String(filepath.Join(dr.Prefix, date)),
	})
	if err != nil {
		return domain.DailyPlayersData{}, fmt.Errorf("s3.NewDailyPlayersDataRepository failed to download from s3: %w", err)
	}

	var s3Players DailyPlayersData

	if err := json.Unmarshal(buf.Bytes(), &s3Players); err != nil {
		return domain.DailyPlayersData{}, fmt.Errorf("s3.NewDailyPlayersDataRepository failed to unmarshal data: %w", err)
	}

	return toDomainDailyPlayersData(s3Players), nil
}

func toDomainDailyPlayersData(data DailyPlayersData) domain.DailyPlayersData {
	domainDailyPlayersData := make(domain.DailyPlayersData)

	for k, v := range data {
		domainDailyPlayersData[k] = domain.Player(v)
	}

	return domainDailyPlayersData
}

func toS3DailyPlayersData(data domain.DailyPlayersData) DailyPlayersData {
	fsDailyPlayersData := make(DailyPlayersData)

	for k, v := range data {
		fsDailyPlayersData[k] = Player(v)
	}

	return fsDailyPlayersData
}
