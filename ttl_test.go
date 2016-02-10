package gon3

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestTurtlePositive(t *testing.T) {

	verbosity := 0

	for _, testName := range positiveParserTests {
		basePath := "./tests/turtle/"
		ttlFile := basePath + testName + ".ttl"
		ntFile := basePath + testName + ".nt"
		ttlBytes, err := ioutil.ReadFile(ttlFile)
		if err != nil {
			t.Fatalf("Error reading file %s", ttlFile)
		}
		ntBytes, err := ioutil.ReadFile(ntFile)
		if err != nil {
			t.Fatalf("Error reading file %s", ntFile)
		}
		if verbosity > 0 {
			fmt.Printf("\nStarting test %s\n", testName)
		}
		ttlGraph, err := NewParser(string(ttlBytes)).Parse()
		if err != nil {
			t.Fatalf("Test %s failed. Error parsing %s\n(%s)", testName, ttlFile, err)
		}
		ntGraph, err := NewParser(string(ntBytes)).Parse()
		if err != nil {
			t.Fatalf("Test %s failed. Error parsing %s\n(%s)", testName, ntFile, err)
		}
		if !ntGraph.IsomorphicTo(ttlGraph) {
			t.Fatalf("Test %s failed. Graphs not isomorphic.\nttl graph:\n%s\nnt graph:\n%s", testName, ttlGraph, ntGraph)
		}
		if verbosity > 0 {
			fmt.Printf("Test %s passed.\n", testName)
		}
	}
}

func TestTurtlePositiveNoIso(t *testing.T) {

	verbosity := 0

	for _, testName := range positiveParserTestsNoIso {
		basePath := "./tests/turtle/"
		ttlFile := basePath + testName + ".ttl"
		ttlBytes, err := ioutil.ReadFile(ttlFile)
		if err != nil {
			t.Fatalf("Error reading file %s", ttlFile)
		}
		if verbosity > 0 {
			fmt.Printf("\nStarting test %s\n", testName)
		}
		_, err = NewParser(string(ttlBytes)).Parse()
		if err != nil {
			t.Fatalf("Test %s failed. Error parsing %s\n(%s)", testName, ttlFile, err)
		}
		if verbosity > 0 {
			fmt.Printf("Test %s passed.\n", testName)
		}
	}
}

func TestTurtleNegative(t *testing.T) {
	verbosity := 0
	for _, testName := range negativeParserTests {
		testFile := "./tests/turtle/" + testName
		b, err := ioutil.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Error reading test file %s", testFile)
		}
		if verbosity > 0 {
			fmt.Printf("\nStarting test %s\n", testName)
		}
		p := NewParser(string(b))
		_, err = p.Parse()
		if err == nil {
			t.Fatalf("Test %s failed: %s", testName, err)
		}
		if verbosity > 0 {
			fmt.Printf("Test %s passed.\n", testName)
		}
	}
}

var negativeParserTests []string = []string{
	// TODO: unicode escaped chars must conform to rules
	//"turtle-eval-bad-01.ttl",
	//"turtle-eval-bad-02.ttl",
	//"turtle-eval-bad-03.ttl",
	"turtle-eval-bad-04.ttl",
	"turtle-syntax-bad-base-01.ttl",
	"turtle-syntax-bad-base-02.ttl",
	"turtle-syntax-bad-base-03.ttl",
	"turtle-syntax-bad-blank-label-dot-end.ttl",
	// TODO: not failing bad string escapes
	//"turtle-syntax-bad-esc-01.ttl",
	// TODO: causing deadlock
	//"turtle-syntax-bad-esc-02.ttl",
	//"turtle-syntax-bad-esc-03.ttl",
	//"turtle-syntax-bad-esc-04.ttl",
	"turtle-syntax-bad-kw-01.ttl",
	"turtle-syntax-bad-kw-02.ttl",
	"turtle-syntax-bad-kw-03.ttl",
	"turtle-syntax-bad-kw-04.ttl",
	"turtle-syntax-bad-kw-05.ttl",
	"turtle-syntax-bad-lang-01.ttl",
	"turtle-syntax-bad-LITERAL2_with_langtag_and_datatype.ttl",
	"turtle-syntax-bad-ln-dash-start.ttl",
	"turtle-syntax-bad-ln-escape-start.ttl",
	"turtle-syntax-bad-ln-escape.ttl",
	"turtle-syntax-bad-missing-ns-dot-end.ttl",
	"turtle-syntax-bad-missing-ns-dot-start.ttl",
	"turtle-syntax-bad-n3-extras-01.ttl",
	"turtle-syntax-bad-n3-extras-02.ttl",
	"turtle-syntax-bad-n3-extras-03.ttl",
	// TODO: hangs
	//"turtle-syntax-bad-n3-extras-04.ttl",
	"turtle-syntax-bad-n3-extras-05.ttl",
	"turtle-syntax-bad-n3-extras-06.ttl",
	"turtle-syntax-bad-n3-extras-07.ttl",
	"turtle-syntax-bad-n3-extras-08.ttl",
	"turtle-syntax-bad-n3-extras-09.ttl",
	"turtle-syntax-bad-n3-extras-10.ttl",
	"turtle-syntax-bad-n3-extras-11.ttl",
	"turtle-syntax-bad-n3-extras-12.ttl",
	"turtle-syntax-bad-n3-extras-13.ttl",
	"turtle-syntax-bad-ns-dot-end.ttl",
	"turtle-syntax-bad-ns-dot-start.ttl",
	"turtle-syntax-bad-num-01.ttl",
	"turtle-syntax-bad-num-02.ttl",
	"turtle-syntax-bad-num-03.ttl",
	"turtle-syntax-bad-num-04.ttl",
	"turtle-syntax-bad-num-05.ttl",
	"turtle-syntax-bad-number-dot-in-anon.ttl",
	"turtle-syntax-bad-pname-01.ttl",
	"turtle-syntax-bad-pname-02.ttl",
	"turtle-syntax-bad-pname-03.ttl",
	"turtle-syntax-bad-prefix-01.ttl",
	"turtle-syntax-bad-prefix-02.ttl",
	"turtle-syntax-bad-prefix-03.ttl",
	"turtle-syntax-bad-prefix-04.ttl",
	"turtle-syntax-bad-prefix-05.ttl",
	"turtle-syntax-bad-string-01.ttl",
	"turtle-syntax-bad-string-02.ttl",
	// TODO: hangs
	//"turtle-syntax-bad-string-03.ttl",
	//"turtle-syntax-bad-string-04.ttl",
	//"turtle-syntax-bad-string-05.ttl",
	"turtle-syntax-bad-string-06.ttl",
	"turtle-syntax-bad-string-07.ttl",
	"turtle-syntax-bad-struct-01.ttl",
	"turtle-syntax-bad-struct-02.ttl",
	"turtle-syntax-bad-struct-03.ttl",
	"turtle-syntax-bad-struct-04.ttl",
	"turtle-syntax-bad-struct-05.ttl",
	"turtle-syntax-bad-struct-06.ttl",
	"turtle-syntax-bad-struct-07.ttl",
	"turtle-syntax-bad-struct-08.ttl",
	"turtle-syntax-bad-struct-09.ttl",
	"turtle-syntax-bad-struct-10.ttl",
	"turtle-syntax-bad-struct-11.ttl",
	"turtle-syntax-bad-struct-12.ttl",
	"turtle-syntax-bad-struct-13.ttl",
	"turtle-syntax-bad-struct-14.ttl",
	"turtle-syntax-bad-struct-15.ttl",
	"turtle-syntax-bad-struct-16.ttl",
	"turtle-syntax-bad-struct-17.ttl",
	"turtle-syntax-bad-uri-01.ttl",
	"turtle-syntax-bad-uri-02.ttl",
	// TODO: causes deadlock
	//"turtle-syntax-bad-uri-03.ttl",
	// TODO: char escapes not allowed in uri
	//"turtle-syntax-bad-uri-04.ttl",
	//"turtle-syntax-bad-uri-05.ttl",
}

var positiveParserTests []string = []string{
	"bareword_a_predicate",
	"bareword_decimal",
	"bareword_double",
	"blankNodePropertyList_as_object",
	"blankNodePropertyList_as_subject",
	"blankNodePropertyList_containing_collection",
	"blankNodePropertyList_with_multiple_triples",
	"collection_object",
	"collection_subject",
	"comment_following_localName",
	"comment_following_PNAME_NS",
	"default_namespace_IRI",
	"double_lower_case_e",
	"empty_collection",
	"first",
	"HYPHEN_MINUS_in_localName",
	"IRIREF_datatype",
	"IRI_subject",
	"IRI_with_all_punctuation",
	"IRI_with_eight_digit_numeric_escape",
	"IRI_with_four_digit_numeric_escape",
	"labeled_blank_node_object",
	"labeled_blank_node_subject",
	"labeled_blank_node_with_leading_digit",
	"labeled_blank_node_with_leading_underscore",
	"labeled_blank_node_with_non_leading_extras",
	"labeled_blank_node_with_PN_CHARS_BASE_character_boundaries",
	"langtagged_LONG",
	"langtagged_LONG_with_subtag",
	"langtagged_non_LONG",
	"lantag_with_subtag",
	"last",
	"LITERAL1_all_controls",
	"LITERAL1_all_punctuation",
	"LITERAL1_ascii_boundaries",
	"LITERAL1",
	"LITERAL1_with_UTF8_boundaries",
	"LITERAL2_ascii_boundaries",
	"LITERAL2",
	"LITERAL2_with_UTF8_boundaries",
	"literal_false",
	"LITERAL_LONG1_ascii_boundaries",
	"LITERAL_LONG1",
	"LITERAL_LONG1_with_1_squote",
	"LITERAL_LONG1_with_2_squotes",
	"LITERAL_LONG1_with_UTF8_boundaries",
	"LITERAL_LONG2_ascii_boundaries",
	"LITERAL_LONG2",
	"LITERAL_LONG2_with_1_squote",
	"LITERAL_LONG2_with_2_squotes",
	"LITERAL_LONG2_with_REVERSE_SOLIDUS",
	"LITERAL_LONG2_with_UTF8_boundaries",
	"literal_true",
	"literal_with_BACKSPACE",
	"literal_with_CARRIAGE_RETURN",
	"literal_with_CHARACTER_TABULATION",
	"literal_with_escaped_BACKSPACE",
	"literal_with_escaped_CARRIAGE_RETURN",
	"literal_with_escaped_CHARACTER_TABULATION",
	"literal_with_escaped_FORM_FEED",
	"literal_with_escaped_LINE_FEED",
	"literal_with_FORM_FEED",
	"literal_with_LINE_FEED",
	"literal_with_numeric_escape4",
	"literal_with_numeric_escape8",
	"literal_with_REVERSE_SOLIDUS",
	"localName_with_assigned_nfc_bmp_PN_CHARS_BASE_character_boundaries",
	"localName_with_assigned_nfc_PN_CHARS_BASE_character_boundaries",
	"localname_with_COLON",
	"localName_with_leading_digit",
	"localName_with_leading_underscore",
	"localName_with_nfc_PN_CHARS_BASE_character_boundaries",
	"localName_with_non_leading_extras",
	// TODO: proper unescaping
	//"manifest.ttl",
	"negative_numeric",
	"nested_blankNodePropertyLists",
	"nested_collection",
	"number_sign_following_localName",
	"number_sign_following_PNAME_NS",
	"numeric_with_leading_0",
	"objectList_with_two_objects",
	"old_style_base",
	"old_style_prefix",
	"percent_escaped_localName",
	"positive_numeric",
	"predicateObjectList_with_two_objectLists",
	"prefixed_IRI_object",
	"prefixed_IRI_predicate",
	"prefixed_name_datatype",
	"prefix_only_IRI",
	"prefix_reassigned_and_used",
	"prefix_with_non_leading_extras",
	"prefix_with_PN_CHARS_BASE_character_boundaries",
	"repeated_semis_at_end",
	"repeated_semis_not_at_end",
	"reserved_escaped_localName",
	"sole_blankNodePropertyList",
	"SPARQL_style_base",
	"SPARQL_style_prefix",
	"turtle-eval-struct-01",
	"turtle-eval-struct-02",
	"turtle-subm-01",
	"turtle-subm-02",
	"turtle-subm-03",
	"turtle-subm-04",
	"turtle-subm-05",
	"turtle-subm-06",
	"turtle-subm-07",
	"turtle-subm-08",
	"turtle-subm-09",
	"turtle-subm-10",
	"turtle-subm-11",
	"turtle-subm-12",
	"turtle-subm-13",
	"turtle-subm-14",
	"turtle-subm-15",
	"turtle-subm-16",
	"turtle-subm-17",
	"turtle-subm-18",
	"turtle-subm-19",
	"turtle-subm-20",
	"turtle-subm-21",
	"turtle-subm-22",
	"turtle-subm-23",
	"turtle-subm-24",
	"turtle-subm-25",
	"turtle-subm-26",
	"turtle-subm-27",
	"two_LITERAL_LONG2s",
	"underscore_in_localName",
}

var positiveParserTestsNoIso []string = []string{
	"bareword_integer",
	"anonymous_blank_node_object",
	"anonymous_blank_node_subject",
	"turtle-syntax-base-01",
	"turtle-syntax-base-02",
	"turtle-syntax-base-03",
	"turtle-syntax-base-04",
	"turtle-syntax-blank-label",
	"turtle-syntax-bnode-01",
	"turtle-syntax-bnode-02",
	"turtle-syntax-bnode-03",
	"turtle-syntax-bnode-04",
	"turtle-syntax-bnode-05",
	"turtle-syntax-bnode-06",
	"turtle-syntax-bnode-07",
	"turtle-syntax-bnode-08",
	"turtle-syntax-bnode-09",
	"turtle-syntax-bnode-10",
	"turtle-syntax-datatypes-01",
	"turtle-syntax-datatypes-02",
	"turtle-syntax-file-01",
	"turtle-syntax-file-02",
	"turtle-syntax-file-03",
	"turtle-syntax-kw-01",
	"turtle-syntax-kw-02",
	"turtle-syntax-kw-03",
	"turtle-syntax-lists-01",
	"turtle-syntax-lists-02",
	"turtle-syntax-lists-03",
	"turtle-syntax-lists-04",
	"turtle-syntax-lists-05",
	"turtle-syntax-ln-colons",
	"turtle-syntax-ln-dots",
	"turtle-syntax-ns-dots",
	"turtle-syntax-number-01",
	"turtle-syntax-number-02",
	"turtle-syntax-number-03",
	"turtle-syntax-number-04",
	"turtle-syntax-number-05",
	"turtle-syntax-number-06",
	"turtle-syntax-number-07",
	"turtle-syntax-number-08",
	"turtle-syntax-number-09",
	"turtle-syntax-number-10",
	"turtle-syntax-number-11",
	"turtle-syntax-pname-esc-01",
	"turtle-syntax-pname-esc-02",
	"turtle-syntax-pname-esc-03",
	"turtle-syntax-prefix-01",
	// this library will NOT allow mixed-case prefix declarations
	//"turtle-syntax-prefix-02",
	"turtle-syntax-prefix-03",
	"turtle-syntax-prefix-04",
	"turtle-syntax-prefix-05",
	"turtle-syntax-prefix-06",
	"turtle-syntax-prefix-07",
	"turtle-syntax-prefix-08",
	"turtle-syntax-prefix-09",
	"turtle-syntax-str-esc-01",
	"turtle-syntax-str-esc-02",
	"turtle-syntax-str-esc-03",
	"turtle-syntax-string-01",
	"turtle-syntax-string-02",
	"turtle-syntax-string-03",
	"turtle-syntax-string-04",
	"turtle-syntax-string-05",
	"turtle-syntax-string-06",
	"turtle-syntax-string-07",
	"turtle-syntax-string-08",
	"turtle-syntax-string-09",
	"turtle-syntax-string-10",
	"turtle-syntax-string-11",
	"turtle-syntax-struct-01",
	"turtle-syntax-struct-02",
	"turtle-syntax-struct-03",
	"turtle-syntax-struct-04",
	"turtle-syntax-struct-05",
	"turtle-syntax-uri-01",
	"turtle-syntax-uri-02",
	"turtle-syntax-uri-03",
	"turtle-syntax-uri-04",
}
