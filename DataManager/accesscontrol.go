package main // github.com/tw/accesscontrol

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const aclRule string = "{srvName:[_a-zA-Z0-9-]+}"

var USE_DEFAULTS bool = true

const DM_HIST_OP string = "/datamanager" + "/historian"
const DM_PROXY_OP string = "/datamanager" + "/proxy"
const DM_HIST_OP_WS string = DM_HIST_OP + "/ws"

const ACL_METHODS string = "gpPd"
const ACL_METHOD_GET string = "g"
const ACL_METHOD_PUT string = "p"
const ACL_METHOD_POST string = "P"
const ACL_METHOD_DEL string = "d"

const METHOD_GET string = "GET"
const METHOD_PUT string = "PUT"
const METHOD_POST string = "POST"
const METHOD_DEL string = "DELETE"

const ACL_LINE_COMMENT1 string = "#"
const ACL_LINE_COMMENT2 string = ";"
const ACL_SYS_WILDCARD string = "$SYS"
const ACL_SRV_WILDCARD string = "*"

const ACL_SYS_DELIM string = ":"
const ACL_RULE_DELIM string = ","
const ACL_PATH_DELIM string = "@"
const ACL_PATH_SEPARATOR string = "/"
const ACL_DEFAULT_RULE string = "$SYS: gpPd@$SYS/*"

type AclEntry struct {
	Operations string
	Path       string
}

type AclRule struct {
	SystemName string
	Acls       []AclEntry
}

var rules []AclRule

func aclInit(aclFile string) error {

	ret := load(aclFile)
	/*for i, rule := range rules {
		fmt.Printf("%d, %+v\n", i, rule)
	}*/

	/*authorized, err := checkRequest("mulle-342", "GET", "/datamanager/historian/ws/mulle-342/temp")
	if err != nil {
		fmt.Printf("Auth check error: %s\n", err)
	}

	if authorized {
		fmt.Printf("System is ok\n")
	} else {
		fmt.Printf("System is evil\n")
	}*/

	return ret
}

func checkRequest(systemCN string, operation string, path string) (bool, error) {
	fmt.Printf("checkRequest(CN: %s, op: %s for '%s')\n", systemCN, operation, path)

	var endPath string = ""
	if strings.Contains(path, DM_HIST_OP_WS) {
		var idx int = strings.Index(path, DM_HIST_OP_WS) + len(DM_HIST_OP_WS) + 1
		if idx < len(path) {
			endPath = path[idx:]
		}
	} else if strings.Contains(path, DM_HIST_OP) {
		var idx int = strings.Index(path, DM_HIST_OP) + len(DM_HIST_OP) + 1
		if idx < len(path) {
			endPath = path[idx:]
		}
	} else if strings.Contains(path, DM_PROXY_OP) {
		var idx int = strings.Index(path, DM_PROXY_OP) + len(DM_PROXY_OP) + 1
		if idx < len(path) {
			endPath = path[idx:]
		}
	}

	//fmt.Printf("endPath: %s\n", endPath)
	if endPath == "" {
		return false, errors.New("Illegal path")
	}

	var targetPath []string = strings.Split(endPath, ACL_PATH_SEPARATOR)
	var op string = ""

	switch strings.TrimSpace(strings.ToUpper(operation)) {
	case METHOD_GET:
		op = ACL_METHOD_GET
	case METHOD_PUT:
		op = ACL_METHOD_PUT
		break
	case METHOD_POST:
		op = ACL_METHOD_POST
		break
	case METHOD_DEL:
		op = ACL_METHOD_DEL
		break
	default:
		return false, errors.New("Illegal operation: " + operation)
	}

	//fmt.Printf("op: %s\nTarget SYS: %s\nTarget SRV: %s\n", op, targetPath[0], targetPath[1])

	for i, rule := range rules {
		if strings.EqualFold(rule.SystemName, systemCN) {
			fmt.Printf("RULE[%v] matches systemName: %+v\n", i, rule)
			for _, acl := range rule.Acls { //j
				//fmt.Printf("Checking [%d] => %+v\n", j, acl)

				var pathParts []string = strings.Split(acl.Path, ACL_PATH_SEPARATOR)
				pathSystem := strings.TrimSpace(pathParts[0])
				pathService := strings.TrimSpace(pathParts[1])

				if acl.Path == endPath {
					if strings.Contains(acl.Operations, op) {
						return true, nil
					}
				} else if pathSystem == targetPath[0] && pathService == ACL_SRV_WILDCARD {
					if strings.Contains(acl.Operations, op) {
						return true, nil
					}
				} else if pathSystem == ACL_SRV_WILDCARD && pathService == ACL_SRV_WILDCARD {
					if strings.Contains(acl.Operations, op) {
						return true, nil
					}
				}
			}

		} else if rule.SystemName == ACL_SYS_WILDCARD {
			for _, acl := range rule.Acls {
				var pathParts []string = strings.Split(acl.Path, ACL_PATH_SEPARATOR)
				pathSystem := ""
				pathService := ""
				if len(pathParts) == 2 {
					pathSystem = strings.TrimSpace(pathParts[0])
					pathService = strings.TrimSpace(pathParts[1])
				}

				if len(targetPath) == 1 {
					if pathSystem == ACL_SYS_WILDCARD && targetPath[0] == systemCN {
						if strings.Contains(acl.Operations, op) {
							return true, nil
						}
					}
				} else {
					if pathSystem == ACL_SYS_WILDCARD && strings.EqualFold(targetPath[0], systemCN) && (pathService == ACL_SRV_WILDCARD || pathService == targetPath[1]) {
						if strings.Contains(acl.Operations, op) {
							return true, nil
						}
					}
				}

			}
		}
	}

	return false, nil
}

func addACLEntries(rulesData []string) ([]AclEntry, error) {
	//log.Printf("addACLEntries:\n")
	var rules []AclEntry = make([]AclEntry, 0)

	for _, rule := range rulesData {
		rule = strings.TrimSpace(rule)
		//fmt.Printf("rulesData[%v]: %s\n", i, rule)

		var parts []string = strings.Split(rule, ACL_PATH_DELIM)
		var operations string = strings.TrimSpace(parts[0]) //.toLowerCase()
		var path string = strings.TrimSpace(parts[1])

		var r AclEntry = AclEntry{}
		r.Operations = operations
		r.Path = path
		rules = append(rules, r)
	}

	return rules, nil
}

/*
func addACLEntry(acls []AclEntry, entryString string) []AclEntry {
	fmt.Printf("addACLEntry(%s):\n", entryString)
	//var ret []AclEntry = make([]AclEntry, 0)
	var parts []string = strings.Split(entryString, ACL_PATH_DELIM)
	var operations string = strings.TrimSpace(parts[0]) //.toLowerCase()
	var path string = strings.TrimSpace(parts[1])

	//for i=0; i< len(operations), i++ {
	//fmt.Printf("rule[%d]: %s\n", i, part)
	//}

	var acl AclEntry = AclEntry{}
	acl.Operations = operations
	acl.Path = path
	acls = append(acls, acl)

	return acls
}
*/

func load(aclFile string) error {
	fmt.Printf("Loading ACL file: %s\n", aclFile)

	fileData, err := ioutil.ReadFile(aclFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rules = make([]AclRule, 0)

	sliceData := strings.Split(string(fileData), "\n")
	for _, line := range sliceData {
		line = strings.TrimSpace(line)
		if !(strings.HasPrefix(line, ACL_LINE_COMMENT1) || strings.HasPrefix(line, ACL_LINE_COMMENT2)) && line != "" {
			//fmt.Printf("LINE(%d): %s\n", i, line)

			s := strings.Split(line, ACL_SYS_DELIM)
			//fmt.Println(s)
			var rule AclRule = AclRule{SystemName: s[0]}
			var rulesData []string = strings.Split(strings.TrimSpace(s[1]), ACL_RULE_DELIM)
			rule.Acls, _ = addACLEntries(rulesData)

			rules = append(rules, rule)
		}
	}

	return nil
}
