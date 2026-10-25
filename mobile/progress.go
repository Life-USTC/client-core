package mobile

import "github.com/Life-USTC/client-core/pagination"

// PageProgress exposes accepted-page progress through gomobile scalars.
type PageProgress struct {
	progress pagination.Progress
}

func NewPageProgress(firstPage int64) *PageProgress {
	return &PageProgress{progress: pagination.NewProgress(int(firstPage))}
}

func (p *PageProgress) NextPage() int64 {
	return int64(p.progress.NextPage())
}

func (p *PageProgress) ReceivedCount() int64 {
	return int64(p.progress.ReceivedCount())
}

// AfterPage returns a new binding value without changing the previous progress.
func (p *PageProgress) AfterPage(itemCount int64, hasMore bool) *PageProgress {
	return &PageProgress{progress: p.progress.AfterPage(int(itemCount), hasMore)}
}
