package psutil

// Get various status from /proc/(pid)/status
//func fillFromStatus(pid int) error {
//	statPath := fmt.Sprintf("/proc/%d/status", pid)
//	contents, err := ioutil.ReadFile(statPath)
//	if err != nil {
//		return err
//	}
//	lines := strings.Split(string(contents), "\n")
//	p.numCtxSwitches = &NumCtxSwitchesStat{}
//	p.memInfo = &MemoryInfoStat{}
//	p.sigInfo = &SignalInfoStat{}
//	for _, line := range lines {
//		tabParts := strings.SplitN(line, "\t", 2)
//		if len(tabParts) < 2 {
//			continue
//		}
//		value := tabParts[1]
//		switch strings.TrimRight(tabParts[0], ":") {
//		case "Name":
//			p.name = strings.Trim(value, " \t")
//			if len(p.name) >= 15 {
//				cmdlineSlice, err := p.CmdlineSlice()
//				if err != nil {
//					return err
//				}
//				if len(cmdlineSlice) > 0 {
//					extendedName := filepath.Base(cmdlineSlice[0])
//					if strings.HasPrefix(extendedName, p.name) {
//						p.name = extendedName
//					} else {
//						p.name = cmdlineSlice[0]
//					}
//				}
//			}
//		case "State":
//			p.status = value[0:1]
//		case "PPid", "Ppid":
//			pval, err := strconv.ParseInt(value, 10, 32)
//			if err != nil {
//				return err
//			}
//			p.parent = int32(pval)
//		case "Tgid":
//			pval, err := strconv.ParseInt(value, 10, 32)
//			if err != nil {
//				return err
//			}
//			p.tgid = int32(pval)
//		case "Uid":
//			p.uids = make([]int32, 0, 4)
//			for _, i := range strings.Split(value, "\t") {
//				v, err := strconv.ParseInt(i, 10, 32)
//				if err != nil {
//					return err
//				}
//				p.uids = append(p.uids, int32(v))
//			}
//		case "Gid":
//			p.gids = make([]int32, 0, 4)
//			for _, i := range strings.Split(value, "\t") {
//				v, err := strconv.ParseInt(i, 10, 32)
//				if err != nil {
//					return err
//				}
//				p.gids = append(p.gids, int32(v))
//			}
//		case "Threads":
//			v, err := strconv.ParseInt(value, 10, 32)
//			if err != nil {
//				return err
//			}
//			p.numThreads = int32(v)
//		case "voluntary_ctxt_switches":
//			v, err := strconv.ParseInt(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.numCtxSwitches.Voluntary = v
//		case "nonvoluntary_ctxt_switches":
//			v, err := strconv.ParseInt(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.numCtxSwitches.Involuntary = v
//		case "VmRSS":
//			value := strings.Trim(value, " kB") // remove last "kB"
//			v, err := strconv.ParseUint(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.memInfo.RSS = v * 1024
//		case "VmSize":
//			value := strings.Trim(value, " kB") // remove last "kB"
//			v, err := strconv.ParseUint(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.memInfo.VMS = v * 1024
//		case "VmSwap":
//			value := strings.Trim(value, " kB") // remove last "kB"
//			v, err := strconv.ParseUint(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.memInfo.Swap = v * 1024
//		case "VmHWM":
//			value := strings.Trim(value, " kB") // remove last "kB"
//			v, err := strconv.ParseUint(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.memInfo.HWM = v * 1024
//		case "VmData":
//			value := strings.Trim(value, " kB") // remove last "kB"
//			v, err := strconv.ParseUint(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.memInfo.Data = v * 1024
//		case "VmStk":
//			value := strings.Trim(value, " kB") // remove last "kB"
//			v, err := strconv.ParseUint(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.memInfo.Stack = v * 1024
//		case "VmLck":
//			value := strings.Trim(value, " kB") // remove last "kB"
//			v, err := strconv.ParseUint(value, 10, 64)
//			if err != nil {
//				return err
//			}
//			p.memInfo.Locked = v * 1024
//		case "SigPnd":
//			v, err := strconv.ParseUint(value, 16, 64)
//			if err != nil {
//				return err
//			}
//			p.sigInfo.PendingThread = v
//		case "ShdPnd":
//			v, err := strconv.ParseUint(value, 16, 64)
//			if err != nil {
//				return err
//			}
//			p.sigInfo.PendingProcess = v
//		case "SigBlk":
//			v, err := strconv.ParseUint(value, 16, 64)
//			if err != nil {
//				return err
//			}
//			p.sigInfo.Blocked = v
//		case "SigIgn":
//			v, err := strconv.ParseUint(value, 16, 64)
//			if err != nil {
//				return err
//			}
//			p.sigInfo.Ignored = v
//		case "SigCgt":
//			v, err := strconv.ParseUint(value, 16, 64)
//			if err != nil {
//				return err
//			}
//			p.sigInfo.Caught = v
//		}
//
//	}
//	return nil
//}
