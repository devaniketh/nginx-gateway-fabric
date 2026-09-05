package telemetry

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("parseSnippetValueIntoDirectives", func() {
	DescribeTable("extracting directive names from snippet values",
		func(input string, expected []string) {
			Expect(parseSnippetValueIntoDirectives(input)).To(Equal(expected))
		},
		Entry("empty snippet", "", nil),
		Entry("whitespace and comment only", "   \n\t  # comment only with ; and {\n   ", nil),
		Entry("simple single directive", "worker_priority 0;", []string{"worker_priority"}),
		Entry("multiple directives on one line", "aio on; keepalive_timeout 65s;", []string{"aio", "keepalive_timeout"}),
		Entry("directive without arguments", "aio; worker_processes auto;", []string{"aio", "worker_processes"}),
		Entry("multi-line directives with comments", "worker_priority 1;\n# comment with ; semicolon\nworker_rlimit_nofile 50;\n", []string{"worker_priority", "worker_rlimit_nofile"}),
		Entry("inline comment with semicolon", "worker_processes auto; # inline; comment\nworker_cpu_affinity auto;", []string{"worker_processes", "worker_cpu_affinity"}),
		Entry("comment before semicolon", "worker_processes auto # comment \n ;", []string{"worker_processes"}),
		Entry("semicolon in double-quoted string", `proxy_set_header hello "myvalue;abc";`, []string{"proxy_set_header"}),
		Entry("semicolon in single-quoted string", `proxy_set_header hello 'myvalue;abc';`, []string{"proxy_set_header"}),
		Entry("escaped quote in double-quoted string", `proxy_set_header foo "bar\"baz;qux"; proxy_set_header test 1;`, []string{"proxy_set_header", "proxy_set_header"}),
		Entry("escaped quote in single-quoted string", `proxy_set_header foo 'bar\'baz;qux'; proxy_set_header test 1;`, []string{"proxy_set_header", "proxy_set_header"}),
		Entry("issue 2702 map block directive", "aio on; map $http_x $x {default 0; 'abc' 1;}", []string{"aio", "map"}),
		Entry("block directive with types", "types { text/html html; image/png png; }", []string{"types"}),
		Entry("block with quoted braces", `map $x $y { "~^foo{" 1; "bar}" 2; default 0; } aio on;`, []string{"map", "aio"}),
		Entry("block with comment containing braces", "map $x $y {\n# { comment }\ndefault 0;\n} aio on;", []string{"map", "aio"}),
		Entry("multiple block directives", "map $a $b { default 0; } map $c $d { default 1; }", []string{"map", "map"}),
		Entry("nested block braces", "block_directive arg { inner { default 0; } } other_dir 1;", []string{"block_directive", "other_dir"}),
	)
})
