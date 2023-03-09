package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pduzinki/fpl-price-checker/internal/config"
	"github.com/pduzinki/fpl-price-checker/internal/domain"
	"github.com/pduzinki/fpl-price-checker/internal/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Player struct {
	ID         int
	Name       string
	Price      int
	SelectedBy string
}

type DailyPlayersData map[int]Player

type DailyPlayersDataRepository struct {
	Bucket string
	Prefix string
	Client *s3.S3
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
		Bucket: awsConfig.Bucket,
		Prefix: prefix,
		Client: s3.New(sess),
	}, nil
}

func (dr *DailyPlayersDataRepository) Add(ctx context.Context, date string, players domain.DailyPlayersData) error {
	s3Players := toS3DailyPlayersData(players)

	if err := domain.ParseDate(date); err != nil {
		return fmt.Errorf("s3.DailyPlayersDataRepository.Add failed to parse date: %w", err)
	}

	jsonPlayers, err := json.Marshal(s3Players)
	if err != nil {
		return fmt.Errorf("s3.NewDailyPlayersDataRepository failed to mashal data: %w", err)
	}

	_, err = dr.Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: &dr.Bucket,
		Key:    aws.String(filepath.Join(dr.Prefix, date)),
	})
	if err == nil {
		return fmt.Errorf("s3.NewDailyPlayersDataRepository.Add failed: %w", storage.ErrDataAlreadyExists)
	}

	_, err = dr.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: &dr.Bucket,
		Key:    aws.String(filepath.Join(dr.Prefix, date)),
		Body:   bytes.NewReader(jsonPlayers),
	})
	if err != nil {
		return fmt.Errorf("s3.NewDailyPlayersDataRepository failed to upload to s3: %w", err)
	}

	return nil
}

func (dr *DailyPlayersDataRepository) GetByDate(ctx context.Context, date string) (domain.DailyPlayersData, error) {
	out, err := dr.Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &dr.Bucket,
		Key:    aws.String(filepath.Join(dr.Prefix, date)),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, fmt.Errorf("s3.DailyPlayersDataRepository.GetByDate failed to download from s3: %w", storage.ErrDataNotFound)
			}
		}

		return nil, fmt.Errorf("s3.DailyPlayersDataRepository.GetByDate failed to download from s3: %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(out.Body)
	if err != nil {
		return nil, fmt.Errorf("s3.DailyPlayersDataRepository.GetByDate failed to read from buffer: %w", err)
	}

	var s3Players DailyPlayersData

	if err := json.Unmarshal(buf.Bytes(), &s3Players); err != nil {
		return nil, fmt.Errorf("s3.DailyPlayersDataRepository failed to unmarshal data: %w", err)
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
