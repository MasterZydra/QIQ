package ini

var allowedDirectives = map[string]int{
	"allow_url_fopen":               INI_SYSTEM,
	"allow_url_include":             INI_SYSTEM,
	"always_populate_raw_post_data": INI_ALL,
	"error_reporting":               INI_ALL,
	"filter.default":                INI_PERDIR,
	"mbstring.encoding_translation": INI_PERDIR,
	"sys_temp_dir":                  INI_SYSTEM,
	// Language Options
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.language-options
	"short_open_tag":      INI_PERDIR,
	"precision":           INI_ALL,
	"serialize_precision": INI_ALL,
	"disable_functions":   INI_SYSTEM,
	"expose_php":          INI_SYSTEM,
	// Data Handling
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.data-handling
	"arg_separator.input":      INI_SYSTEM,
	"arg_separator.output":     INI_ALL,
	"variables_order":          INI_PERDIR,
	"auto_globals_jit":         INI_PERDIR,
	"register_argc_argv":       INI_PERDIR,
	"enable_post_data_reading": INI_PERDIR,
	"post_max_size":            INI_PERDIR,
	"auto_prepend_file":        INI_PERDIR,
	"default_charset":          INI_ALL,
	"input_encoding":           INI_ALL,
	"output_encoding":          INI_ALL,
	"internal_encoding":        INI_ALL,
	// File Uploads
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.file-uploads
	"file_uploads":            INI_SYSTEM,
	"upload_tmp_dir":          INI_SYSTEM,
	"max_input_nesting_level": INI_PERDIR,
	"max_input_vars":          INI_PERDIR,
	"upload_max_filesize":     INI_PERDIR,
	"max_file_uploads":        INI_PERDIR,
	// Installation / Configuration
	// Spec: https://www.php.net/manual/en/info.configuration.php
	"assert.exception":            INI_ALL,
	"max_execution_time":          INI_ALL,
	"max_input_time":              INI_PERDIR,
	"zend.enable_gc":              INI_ALL,
	"zend.max_allowed_stack_size": INI_SYSTEM,
	"zend.reserved_stack_size":    INI_SYSTEM,
	"fiber.stack_size":            INI_ALL,
	// Paths and Directories
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.path-directory
	"include_path": INI_ALL,
	"open_basedir": INI_ALL,
	// Resource Limits
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.resource-limits
	"max_memory_limit": INI_SYSTEM,
	"memory_limit":     INI_ALL,
	// Session
	// Spec: https://www.php.net/manual/en/session.configuration.php
	"session.save_path":                INI_ALL,
	"session.name":                     INI_ALL,
	"session.save_handler":             INI_ALL,
	"session.auto_start":               INI_PERDIR,
	"session.gc_probability":           INI_ALL,
	"session.gc_divisor":               INI_ALL,
	"session.gc_maxlifetime":           INI_ALL,
	"session.serialize_handler":        INI_ALL,
	"session.cookie_lifetime":          INI_ALL,
	"session.cookie_path":              INI_ALL,
	"session.cookie_domain":            INI_ALL,
	"session.cookie_secure":            INI_ALL,
	"session.cookie_httponly":          INI_ALL,
	"session.cookie_samesite":          INI_ALL,
	"session.use_strict_mode":          INI_ALL,
	"session.use_cookies":              INI_ALL,
	"session.use_only_cookies":         INI_ALL,
	"session.referer_check":            INI_ALL,
	"session.cache_limiter":            INI_ALL,
	"session.cache_expire":             INI_ALL,
	"session.use_trans_sid":            INI_ALL,
	"session.trans_sid_tags":           INI_ALL,
	"session.trans_sid_hosts":          INI_ALL,
	"session.sid_length":               INI_ALL,
	"session.sid_bits_per_character":   INI_ALL,
	"session.upload_progress.enabled":  INI_PERDIR,
	"session.upload_progress.cleanup":  INI_PERDIR,
	"session.upload_progress.prefix":   INI_PERDIR,
	"session.upload_progress.name":     INI_PERDIR,
	"session.upload_progress.freq":     INI_PERDIR,
	"session.upload_progress.min_freq": INI_PERDIR,
	"session.lazy_write":               INI_ALL,
	// Zlib
	// Spec: https://www.php.net/manual/en/zlib.configuration.php
	"zlib.output_compression":       INI_ALL,
	"zlib.output_compression_level": INI_ALL,
	"zlib.output_handler":           INI_ALL,
	// QIQ
	"qiq.case_sensitive_include": INI_SYSTEM,
}

var defaultValues = map[string]string{
	// Ini Directives:
	"allow_url_fopen":               "1",
	"allow_url_include":             "0",
	"always_populate_raw_post_data": "",
	"error_reporting":               "0",
	"filter.default":                "unsafe_raw",
	"mbstring.encoding_translation": "",
	"sys_temp_dir":                  "",
	// Category: Language Options
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.language-options
	"short_open_tag":      "",
	"precision":           "14",
	"serialize_precision": "-1",
	"disable_functions":   "",
	"expose_php":          "",
	// Category: Data Handling
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.data-handling
	"arg_separator.input":      "&",
	"arg_separator.output":     "&",
	"variables_order":          "EGPCS",
	"auto_globals_jit":         "1",
	"register_argc_argv":       "",
	"enable_post_data_reading": "1",
	"post_max_size":            "8M",
	"auto_prepend_file":        "",
	"default_charset":          "UTF-8",
	"input_encoding":           "",
	"output_encoding":          "",
	"internal_encoding":        "",
	// Category: File Uploads
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.file-uploads
	"file_uploads":            "1",
	"upload_tmp_dir":          "",
	"max_input_nesting_level": "64",
	"max_input_vars":          "1000",
	"upload_max_filesize":     "2M",
	"max_file_uploads":        "20",
	// Category: Installation / Configuration
	// Spec: https://www.php.net/manual/en/info.configuration.php
	"assert.exception":            "1",
	"max_execution_time":          "30",
	"max_input_time":              "-1",
	"zend.enable_gc":              "1",
	"zend.max_allowed_stack_size": "0",
	"zend.reserved_stack_size":    "0",
	"fiber.stack_size":            "2M",
	// Category: Paths and Directories
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.path-directory
	"include_path": ".",
	"open_basedir": "",
	// Category: Resource Limits
	// Spec: https://www.php.net/manual/en/ini.core.php#ini.sect.resource-limits
	"max_memory_limit": "-1",
	"memory_limit":     "128M",
	// Category: Session
	// Spec: https://www.php.net/manual/en/session.configuration.php
	"session.save_path":                "",
	"session.name":                     "PHPSESSID",
	"session.save_handler":             "files",
	"session.auto_start":               "0",
	"session.gc_probability":           "1",
	"session.gc_divisor":               "100",
	"session.gc_maxlifetime":           "1440",
	"session.serialize_handler":        "php",
	"session.cookie_lifetime":          "0",
	"session.cookie_path":              "/",
	"session.cookie_domain":            "",
	"session.cookie_secure":            "0",
	"session.cookie_httponly":          "0",
	"session.cookie_samesite":          "",
	"session.use_strict_mode":          "",
	"session.use_cookies":              "1",
	"session.use_only_cookies":         "1",
	"session.referer_check":            "",
	"session.cache_limiter":            "nocache",
	"session.cache_expire":             "180",
	"session.use_trans_sid":            "0",
	"session.trans_sid_tags":           "a=href,area=href,frame=src,form=",
	"session.trans_sid_hosts":          "$_SERVER['HTTP_HOST']",
	"session.sid_length":               "32",
	"session.sid_bits_per_character":   "4",
	"session.upload_progress.enabled":  "1",
	"session.upload_progress.cleanup":  "1",
	"session.upload_progress.prefix":   "upload_progress_",
	"session.upload_progress.name":     "PHP_SESSION_UPLOAD_PROGRESS",
	"session.upload_progress.freq":     "1%",
	"session.upload_progress.min_freq": "1",
	"session.lazy_write":               "1",
	// Category: Zlib
	// Spec: https://www.php.net/manual/en/zlib.configuration.php
	"zlib.output_compression":       "0",
	"zlib.output_compression_level": "-1",
	"zlib.output_handler":           "",
	// Category: QIQ
	"qiq.case_sensitive_include": "0",
}

var boolDirectives = []string{
	"allow_url_fopen",
	"allow_url_include",
	"assert.exception",
	"auto_globals_jit",
	"enable_post_data_reading",
	"expose_php",
	"file_uploads",
	"mbstring.encoding_translation",
	"register_argc_argv",
	"session.auto_start",
	"session.cookie_httponly",
	"session.cookie_secure",
	"session.lazy_write",
	"session.upload_progress.cleanup",
	"session.upload_progress.enabled",
	"session.use_cookies",
	"session.use_only_cookies",
	"session.use_strict_mode",
	"session.use_trans_sid",
	"short_open_tag",
	"zend.enable_gc",
	"zlib.output_compression",
	// QIQ
	"qiq.case_sensitive_include",
}

var intDirectives = []string{
	"error_reporting",
	"fiber.stack_size",
	"max_execution_time",
	"max_file_uploads",
	"max_input_nesting_level",
	"max_input_time",
	"max_input_vars",
	// "max_memory_limit", // TODO Fix ini logic because this can also be "128M"
	// "memory_limit", // TODO Fix ini logic because this can also be "128M"
	"precision",
	"serialize_precision",
	"session.cache_expire",
	"session.cookie_lifetime",
	"session.entropy_length",
	"session.gc_divisor",
	"session.gc_maxlifetime",
	"session.gc_probability",
	"session.sid_bits_per_character",
	"session.sid_length",
	"session.upload_progress.min_freq",
	"zend.max_allowed_stack_size",
	"zend.reserved_stack_size",
	"zlib.output_compression_level",
	"zlib.output_compression",
}
