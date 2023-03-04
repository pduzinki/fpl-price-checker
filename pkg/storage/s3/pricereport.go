package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pduzinki/fpl-price-checker/pkg/config"
	"github.com/pduzinki/fpl-price-checker/pkg/domain"
	"github.com/pduzinki/fpl-price-checker/pkg/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Record struct {
	Name        string
	OldPrice    string
	NewPrice    string
	Description string
}

type PriceChangeReport struct {
	Date    string
	Records []Record
}

type PriceReportRepository struct {
	Bucket string
	Prefix string
	Client *s3.S3
}

func NewPriceReportRepository(awsConfig config.AWSConfig, prefix string) (*PriceReportRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           &awsConfig.Region,
		Credentials:      credentials.NewStaticCredentials(awsConfig.ID, awsConfig.Secret, awsConfig.Token),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         &awsConfig.Endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("s3.NewPriceReportRepository failed: %w", err)
	}

	return &PriceReportRepository{
		Bucket: awsConfig.Bucket,
		Prefix: prefix,
		Client: s3.New(sess),
	}, nil
}

func (pr *PriceReportRepository) Add(ctx context.Context, date string, report domain.PriceChangeReport) error {
	s3Report := toS3Report(report)

	jsonReport, err := json.Marshal(s3Report)
	if err != nil {
		return fmt.Errorf("s3.PriceReportRepository.Add failed to marshal data: %w", err)
	}

	// TODO check if object exists
	// _, err = pr.GetByDate(ctx, date)
	// if err == nil {
	// 	return fmt.Errorf("s3.PriceReportRepository.Add failed: %w", storage.ErrDataAlreadyExists)
	// }

	_, err = pr.Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: &pr.Bucket,
		Key:    aws.String(filepath.Join(pr.Prefix, date)),
	})
	if err == nil {
		return fmt.Errorf("s3.PriceReportRepository.Add failed: %w", storage.ErrDataAlreadyExists)
	}

	_, err = pr.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: &pr.Bucket,
		Key:    aws.String(filepath.Join(pr.Prefix, date)),
		Body:   bytes.NewReader(jsonReport),
	})
	if err != nil {
		return fmt.Errorf("s3.PriceReportRepository.Add failed to upload to s3: %w", err)
	}

	return nil
}

func (pr *PriceReportRepository) GetByDate(ctx context.Context, date string) (domain.PriceChangeReport, error) {
	out, err := pr.Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &pr.Bucket,
		Key:    aws.String(filepath.Join(pr.Prefix, date)),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return domain.PriceChangeReport{}, fmt.Errorf("s3.PriceReportRepository.GetByDate failed to download from s3: %w", storage.ErrDataNotFound)
			}
		}

		return domain.PriceChangeReport{}, fmt.Errorf("s3.PriceReportRepository.GetByDate failed to download from s3: %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(out.Body)
	if err != nil {
		return domain.PriceChangeReport{}, fmt.Errorf("s3.PriceReportRepository.GetByDate failed to read from buffer: %w", err)
	}

	var s3Report PriceChangeReport

	if err := json.Unmarshal(buf.Bytes(), &s3Report); err != nil {
		return domain.PriceChangeReport{}, fmt.Errorf("s3.PriceReportRepository.GetByDate failed to unmarshal data: %w", err)
	}

	return toDomainReport(s3Report), nil
}

func toS3Report(report domain.PriceChangeReport) PriceChangeReport {
	records := make([]Record, 0, len(report.Records))

	for _, r := range report.Records {
		records = append(records, Record(r))
	}

	return PriceChangeReport{
		Date:    report.Date,
		Records: records,
	}
}

func toDomainReport(report PriceChangeReport) domain.PriceChangeReport {
	records := make([]domain.Record, 0, len(report.Records))

	for _, r := range report.Records {
		records = append(records, domain.Record(r))
	}

	return domain.PriceChangeReport{
		Date:    report.Date,
		Records: records,
	}
}
