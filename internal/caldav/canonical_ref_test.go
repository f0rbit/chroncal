package caldav

import "testing"

// TestCanonicalObjectRefRejectsCollectionSelfHref guards against a real
// iCloud interop bug: sync-collection sometimes lists the calendar
// collection itself as a "changed" entry, with an href that is missing the
// trailing slash a collection href normally carries. Before this fix, that
// href passed CanonicalObjectRef's trailing-slash check and slipped into a
// calendar-multiget REPORT as a malformed resource reference. iCloud's
// server never answered that request, and chroncal's HTTP client stalled
// for the full 30s timeout on every pull.
func TestCanonicalObjectRefRejectsCollectionSelfHref(t *testing.T) {
	t.Parallel()

	client, err := NewClient(putTestHTTPClient{}, "https://caldav.icloud.com/12345/calendars/AAAA-BBBB/")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// The collection href without its trailing slash, exactly as iCloud
	// returns it inside a sync-collection change entry.
	_, err = client.CanonicalObjectRef(
		"/12345/calendars/AAAA-BBBB/",
		"/12345/calendars/AAAA-BBBB",
	)
	if err == nil {
		t.Fatal("CanonicalObjectRef: want error for a collection self-href, got nil")
	}
}

func TestCanonicalObjectRefRejectsCollectionSelfHrefWithSlash(t *testing.T) {
	t.Parallel()

	client, err := NewClient(putTestHTTPClient{}, "https://caldav.icloud.com/12345/calendars/AAAA-BBBB/")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	_, err = client.CanonicalObjectRef(
		"/12345/calendars/AAAA-BBBB/",
		"/12345/calendars/AAAA-BBBB/",
	)
	if err == nil {
		t.Fatal("CanonicalObjectRef: want error for a collection href, got nil")
	}
}

func TestCanonicalObjectRefAcceptsRealResource(t *testing.T) {
	t.Parallel()

	client, err := NewClient(putTestHTTPClient{}, "https://caldav.icloud.com/12345/calendars/AAAA-BBBB/")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	got, err := client.CanonicalObjectRef(
		"/12345/calendars/AAAA-BBBB/",
		"/12345/calendars/AAAA-BBBB/event.ics",
	)
	if err != nil {
		t.Fatalf("CanonicalObjectRef: %v", err)
	}
	want := "/12345/calendars/AAAA-BBBB/event.ics"
	if got != want {
		t.Fatalf("CanonicalObjectRef = %q, want %q", got, want)
	}
}
