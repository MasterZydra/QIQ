package ini

import (
	"QIQ/cmd/qiq/common"
	"fmt"
	"slices"
	"strings"
)

type Writer struct {
	filename string
}

func NewWriter(filename string) *Writer { return &Writer{filename: filename} }

func (w *Writer) getDefaultValue(directive string) (string, bool) {
	value, found := defaultValues[directive]
	if !found {
		return "", found
	}
	if slices.Contains(boolDirectives, directive) {
		if value == "1" {
			return "on", true
		}
		return "off", true
	}
	return value, true
}

func (w *Writer) addDirectiveAndValue(builder *strings.Builder, directive string) error {
	value, found := w.getDefaultValue(directive)
	if !found {
		return fmt.Errorf("failed to get default value for ini directive %s", directive)
	}
	builder.WriteString(directive)
	builder.WriteString(" = ")

	if slices.Contains(boolDirectives, directive) || slices.Contains(intDirectives, directive) {
		builder.WriteString(value)
	} else if value != "" {
		builder.WriteString(`"`)
		builder.WriteString(value)
		builder.WriteString(`"`)
	}
	builder.WriteString("\n")
	return nil
}

func (w *Writer) addHeading(builder *strings.Builder, heading string) {
	builder.WriteString("\n; ----------------------\n; ")
	builder.WriteString(heading)
	builder.WriteString("\n; ----------------------\n\n")
}

func (w *Writer) Write() error {
	var builder strings.Builder

	builder.WriteString("[PHP]\n\n")

	if err := w.addDirectiveAndValue(&builder, "allow_url_fopen"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "allow_url_include"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "always_populate_raw_post_data"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "error_reporting"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "filter.default"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "mbstring.encoding_translation"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "sys_temp_dir"); err != nil {
		return err
	}

	w.addHeading(&builder, "QIQ")

	builder.WriteString("; Enforces case-sensitive behavior for `include` and `require` statements even on Windows systems.\n")
	if err := w.addDirectiveAndValue(&builder, "qiq.case_sensitive_include"); err != nil {
		return err
	}
	builder.WriteString("\n; Enforces strict comparison semantics throughout your codebase.\n; When enabled, the `==` operator behaves like `===`, and both `!=` and `<>` behave like `!==`.\n")
	if err := w.addDirectiveAndValue(&builder, "qiq.strict_comparison"); err != nil {
		return err
	}

	w.addHeading(&builder, "Language Options")

	if err := w.addDirectiveAndValue(&builder, "short_open_tag"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "precision"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "serialize_precision"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "disable_functions"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "expose_php"); err != nil {
		return err
	}

	w.addHeading(&builder, "Resource Limits")

	if err := w.addDirectiveAndValue(&builder, "max_memory_limit"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "memory_limit"); err != nil {
		return err
	}

	w.addHeading(&builder, "Data Handling")

	if err := w.addDirectiveAndValue(&builder, "arg_separator.input"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "arg_separator.output"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "variables_order"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "auto_globals_jit"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "register_argc_argv"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "enable_post_data_reading"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "post_max_size"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "auto_prepend_file"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "default_charset"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "input_encoding"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "output_encoding"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "internal_encoding"); err != nil {
		return err
	}

	w.addHeading(&builder, "File Uploads")

	if err := w.addDirectiveAndValue(&builder, "file_uploads"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "upload_tmp_dir"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "max_input_nesting_level"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "max_input_vars"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "upload_max_filesize"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "max_file_uploads"); err != nil {
		return err
	}

	w.addHeading(&builder, "Installation / Configuration")

	if err := w.addDirectiveAndValue(&builder, "assert.exception"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "max_execution_time"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "max_input_time"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "zend.enable_gc"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "zend.max_allowed_stack_size"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "zend.reserved_stack_size"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "fiber.stack_size"); err != nil {
		return err
	}

	w.addHeading(&builder, "Paths and Directories")

	if err := w.addDirectiveAndValue(&builder, "include_path"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "open_basedir"); err != nil {
		return err
	}

	w.addHeading(&builder, "Zlib")

	if err := w.addDirectiveAndValue(&builder, "zlib.output_compression"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "zlib.output_compression_level"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "zlib.output_handler"); err != nil {
		return err
	}

	w.addHeading(&builder, "Session")

	if err := w.addDirectiveAndValue(&builder, "session.save_path"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.name"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.save_handler"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.auto_start"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.gc_probability"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.gc_divisor"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.gc_maxlifetime"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.serialize_handler"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cookie_lifetime"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cookie_path"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cookie_domain"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cookie_secure"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cookie_httponly"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cookie_samesite"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.use_strict_mode"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.use_cookies"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.use_only_cookies"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.referer_check"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cache_limiter"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.cache_expire"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.use_trans_sid"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.trans_sid_tags"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.trans_sid_hosts"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.sid_length"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.sid_bits_per_character"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.upload_progress.enabled"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.upload_progress.cleanup"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.upload_progress.prefix"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.upload_progress.name"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.upload_progress.freq"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.upload_progress.min_freq"); err != nil {
		return err
	}
	if err := w.addDirectiveAndValue(&builder, "session.lazy_write"); err != nil {
		return err
	}

	return common.WriteFile(w.filename, builder.String())
}
