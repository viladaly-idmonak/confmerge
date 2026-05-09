package merger

import "fmt"

// MergeWithAudit performs a deep merge of src into dst, recording every
// overwrite and addition into the provided AuditLog.
func MergeWithAudit(dst, src map[string]interface{}, source string, log *AuditLog) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range dst {
		result[k] = v
	}
	mergeIntoAudit(result, src, "", source, log)
	return result
}

func mergeIntoAudit(dst, src map[string]interface{}, prefix, source string, log *AuditLog) {
	for k, srcVal := range src {
		path := k
		if prefix != "" {
			path = fmt.Sprintf("%s.%s", prefix, k)
		}

		dstVal, exists := dst[k]
		if !exists {
			dst[k] = srcVal
			log.Record("add", path, nil, srcVal, source)
			continue
		}

		srcMap, srcIsMap := toMap(srcVal)
		dstMap, dstIsMap := toMap(dstVal)

		if srcIsMap && dstIsMap {
			mergeIntoAudit(dstMap, srcMap, path, source, log)
			dst[k] = dstMap
		} else {
			log.Record("override", path, dstVal, srcVal, source)
			dst[k] = srcVal
		}
	}
}
