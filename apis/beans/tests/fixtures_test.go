package gobeansack_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	pgvector "github.com/pgvector/pgvector-go"
	"github.com/soumitsalman/cafecito-api-platform/apis/shared"
)

var (
	fixtureSourceTech  = uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa1")
	fixtureSourceSlash = uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaa2")
	fixturePostID      = uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa10")
	fixtureArticleIDs  = []uuid.UUID{
		uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa01"),
		uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa02"),
		uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa03"),
		uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa04"),
		uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa05"),
		uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa06"),
		uuid.MustParse("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaa07"),
		fixturePostID,
	}
)

func TestMain(m *testing.M) {
	loadTestEnv()
	if err := seedCIFixtures(); err != nil {
		fmt.Fprintf(os.Stderr, "seed CI fixtures: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func loadTestEnv() {
	_ = godotenv.Load()
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		_ = godotenv.Load(filepath.Join(dir, ".env"))
		_ = godotenv.Load(filepath.Join(dir, "..", ".env"))
	}
}

func seedCIFixtures() error {
	conn_str := os.Getenv("PG_CONNECTION_STRING")
	if conn_str == "" {
		return fmt.Errorf("PG_CONNECTION_STRING is required for beans tests")
	}
	ctx := context.Background()
	pool, err := shared.NewConnection(ctx, conn_str)
	if err != nil {
		return err
	}
	defer pool.Close()

	now := time.Now().UTC()
	embedding := pgvector.NewVector(test_query_embedding)

	if _, err := pool.Exec(ctx, `
		INSERT INTO publishers (id, domain_name, base_url, site_name, description, favicon, rss_feed)
		VALUES
			($1, $2, $3, $4, $5, $6, $7),
			($8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (id) DO UPDATE SET
			domain_name = EXCLUDED.domain_name,
			base_url = EXCLUDED.base_url,
			site_name = EXCLUDED.site_name
	`, fixtureSourceTech, "techcrunch.com", "https://techcrunch.com", "TechCrunch",
		"CI fixture publisher", "https://techcrunch.com/favicon.ico", "https://techcrunch.com/feed",
		fixtureSourceSlash, "slashgear.com", "https://slashgear.com", "SlashGear",
		"CI fixture publisher", "https://slashgear.com/favicon.ico", "https://slashgear.com/feed"); err != nil {
		return fmt.Errorf("publishers: %w", err)
	}

	kinds := []string{"news", "news", "news", "news", "news", "news", "blog", "post"}
	sources := []uuid.UUID{fixtureSourceTech, fixtureSourceTech, fixtureSourceSlash, fixtureSourceTech, fixtureSourceTech, fixtureSourceTech, fixtureSourceTech, fixtureSourceTech}
	urls := []string{
		test_article_urls[0], test_article_urls[1], test_article_urls[2], test_article_urls[3], test_article_urls[4],
		"https://techcrunch.com/ci/beans-article-6",
		"https://techcrunch.com/ci/beans-article-7",
		"https://techcrunch.com/ci/beans-post-8",
	}
	for i, id := range fixtureArticleIDs {
		created := now.Add(-time.Duration(i) * time.Hour)
		if err := upsertBean(ctx, pool, id, urls[i], kinds[i], sources[i], created, embedding); err != nil {
			return err
		}
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO related_beans (url, related_url, collected)
		VALUES ($1, $2, NOW()), ($2, $3, NOW())
		ON CONFLICT DO NOTHING
	`, urls[0], urls[1], urls[2]); err != nil {
		return fmt.Errorf("related_beans: %w", err)
	}

	for i := 0; i < 6; i++ {
		if _, err := pool.Exec(ctx, `
			INSERT INTO chatters (chatter_url, url, bean_id, platform, forum, collected, likes, comments, subscribers, shares)
			VALUES ($1, $2, $3, 'reddit', 'r/technology', NOW(), 20, 8, 1000, 2)
		`, fmt.Sprintf("https://reddit.com/r/technology/ci-beans-%d-%s", i+1, fixtureArticleIDs[i].String()[:8]), urls[i], fixtureArticleIDs[i]); err != nil {
			if _, err2 := pool.Exec(ctx, `
				INSERT INTO chatters (chatter_url, url, platform, forum, collected, likes, comments, subscribers, shares)
				VALUES ($1, $2, 'reddit', 'r/technology', NOW(), 20, 8, 1000, 2)
			`, fmt.Sprintf("https://reddit.com/r/technology/ci-beans-%d-%s", i+1, fixtureArticleIDs[i].String()[:8]), urls[i]); err2 != nil {
				return fmt.Errorf("chatters: %w / %w", err, err2)
			}
		}
	}

	if _, err := pool.Exec(ctx, `REFRESH MATERIALIZED VIEW trend_aggregates`); err != nil {
		return fmt.Errorf("refresh trend_aggregates: %w", err)
	}
	return nil
}

func upsertBean(ctx context.Context, pool *pgxpool.Pool, id uuid.UUID, url, kind string, source_id uuid.UUID, created time.Time, embedding pgvector.Vector) error {
	base_url := "https://techcrunch.com"
	if source_id == fixtureSourceSlash {
		base_url = "https://slashgear.com"
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO beans (
			id, url, kind, author, source_id, base_url, image_url, created, collected,
			title, summary, content, restricted_content, embedding,
			categories, sentiments, regions, entities
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $8,
			$9, $10, $11, false, $12,
			$13, $14, $15, $16
		)
		ON CONFLICT (id) DO UPDATE SET
			url = EXCLUDED.url,
			kind = EXCLUDED.kind,
			created = EXCLUDED.created,
			embedding = EXCLUDED.embedding,
			categories = EXCLUDED.categories,
			sentiments = EXCLUDED.sentiments,
			regions = EXCLUDED.regions,
			entities = EXCLUDED.entities
	`, id, url, kind, test_authors[0], source_id, base_url, base_url+"/image.png", created,
		"CI fixture article "+id.String()[:8],
		"CI fixture summary",
		"CI fixture full content",
		embedding,
		[]string{test_categories[0], test_tags[0]},
		[]string{test_sentiments[0]},
		[]string{test_regions[2]},
		[]string{test_entities[0]},
	)
	if err != nil {
		return fmt.Errorf("beans %s: %w", id, err)
	}
	return nil
}
