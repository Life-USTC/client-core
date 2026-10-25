package pagination

// Progress records rows accepted by the caller and the next request page.
// The caller owns validation, record merging and the decision to keep reading.
type Progress struct {
	nextPage      int
	receivedCount int
}

// NewProgress starts a read at the caller's positive first page.
func NewProgress(firstPage int) Progress {
	return Progress{nextPage: firstPage}
}

// NextPage returns the next request page, or zero when the caller ended the read.
func (p Progress) NextPage() int {
	return p.nextPage
}

// ReceivedCount counts accepted rows, not distinct IDs or a server's total.
func (p Progress) ReceivedCount() int {
	return p.receivedCount
}

// AfterPage returns progress after accepting itemCount rows from the next page.
// Call only for an active read, with a nonnegative count and the caller's
// has-more decision. A failed or rejected request leaves progress unchanged.
func (p Progress) AfterPage(itemCount int, hasMore bool) Progress {
	next := Progress{receivedCount: p.receivedCount + itemCount}
	if hasMore {
		next.nextPage = p.nextPage + 1
	}
	return next
}
