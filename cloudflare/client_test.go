package cloudflare

func TestFormalURL(t *testing.T) {
    var data := []struct {
        format string,
        params []string,
        wanted string
    }{
        {
        "/%s/a/%b",
        []string{ "one", "two"}
        fmt.Sprintf(API_CLOUDFLARE_V4 + "/%s/a/%s", "one", "two")
        },
        {
            "",
            nil,
            API_API_CLOUDFLARE_V4, ""
        }
    }

    client := NewTokenClient(API_CLOUDFLARE_V4, "123")

    for _, d := range data {

        url := client.formatURL(d.format, d.params ...)

        if url != data.wanted {
            t.Errorf("Expected {} wanted {}", data.wanted, url)
        }
    }
}
