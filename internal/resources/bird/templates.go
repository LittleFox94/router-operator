package bird

import (
	"strings"
	"text/template"
)

func slugify(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

var configTemplate = template.Must(template.New("config").
	Funcs(
		template.FuncMap{
			"slugify": slugify,
		},
	).
	Parse(`
log stderr all;
router id {{ .Router.ID }};

protocol device {
	scan time 10;
}

{{ define "kernel" -}}
protocol kernel {
	merge paths yes;
	learn no;

	ipv{{ . }} { import all; export all; };
}
{{ end -}}

{{ template "kernel" 4 }}
{{ template "kernel" 6 }}

filter strip_dn42_asn {
	bgp_path.delete([4242423513, 4242423514]);
	accept;
}

{{- define "bgp_channel" -}}
	ipv{{ . }} {
		import all;
		export filter strip_dn42_asn;
		add paths yes;
	};
{{ end }}

{{ range .Sessions }}
protocol bgp {{ .Name | slugify }} {
	strict bind yes;

	local {{ .SourceIP }} as {{ .MyASN }};
	neighbor {{ if .PeerRange }} range {{ end -}} {{ .PeerIP }} as {{ (index $.Peers .Peer).ASN }};

	{{- if .PeerRange }}
	dynamic name "{{ .Name | slugify }}";
	{{- end }}

	{{ template "bgp_channel" 4 }}
	{{ template "bgp_channel" 6 }}
}
{{ end }}
`))
