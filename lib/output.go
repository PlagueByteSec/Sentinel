package lib

import (
	"Sentinel/lib/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func WriteJSON(jsonFileName string) error {
	/*
		Write the summary in JSON format to a file. The default
		directory (output) is used if no custom path is specified with the -j flag.
	*/
	bytes, err := json.MarshalIndent(utils.GJsonResult.Subdomains, "", "	")
	if err != nil {
		utils.Glogger.Println(err)
		return err
	}
	if err := os.WriteFile(jsonFileName, bytes, utils.DefaultPermission); err != nil {
		utils.Glogger.Println(err)
		return errors.New("failed to write JSON to: " + jsonFileName)
	}
	return nil
}

func OutputHandler(streams *utils.FileStreams, client *http.Client, args *utils.Args, params utils.Params) {
	utils.GSubdomBase = utils.SubdomainBase{}
	if args.HttpCode || args.AnalyzeHeader {
		time.Sleep(time.Duration(args.HttpRequestDelay) * time.Millisecond)
	}
	utils.GStdout.Flush()
	/*
		Perform a DNS lookup to determine the IP addresses (IPv4 and IPv6). The addresses will
		be returned as a slice and separated as strings.
	*/
	ipAddrsOut, ipAddrs := IpResolveWrapper(utils.GDnsResolver, args, params)
	if ipAddrs == nil {
		return
	}
	var (
		// Use strings.Builder for better output control
		consoleOutput strings.Builder
		codeFilterExc []string
		codeFilter    []string
	)
	utils.GSubdomBase.Subdomain = append(utils.GSubdomBase.Subdomain, params.Subdomain)
	consoleOutput.WriteString(fmt.Sprintf(" ══[ %s", params.Subdomain))
	/*
		Split the arguments specified by the -f and -e flags by comma.
		The values within the slices will be used to filter the results.
	*/
	delim := ","
	if !strings.Contains(args.ExcHttpCodes, delim) {
		codeFilterExc = []string{args.ExcHttpCodes}
	} else {
		codeFilterExc = strings.Split(args.ExcHttpCodes, delim)
	}
	if !strings.Contains(args.FilHttpCodes, delim) {
		codeFilter = []string{args.FilHttpCodes}
	} else {
		codeFilter = strings.Split(args.FilHttpCodes, delim)
	}
	utils.ResetSlice(&codeFilter)
	utils.ResetSlice(&codeFilterExc)
	if args.HttpCode {
		url := fmt.Sprintf("http://%s", params.Subdomain)
		httpStatusCode := HttpStatusCode(client, url)
		statusCodeConv := strconv.Itoa(httpStatusCode)
		if httpStatusCode == -1 {
			statusCodeConv = utils.NotAvailable
		}
		/*
			Ensure that the status codes are correctly filtered by comparing the
			results with codeFilter and CodeFilterExc.
		*/
		if len(codeFilter) >= 1 && !utils.InArgList(statusCodeConv, codeFilter) ||
			len(codeFilterExc) >= 1 && utils.InArgList(statusCodeConv, codeFilterExc) {
			return
		} else {
			OutputWrapper(ipAddrs, params, streams)
		}
		consoleOutput.WriteString(fmt.Sprintf(", HTTP Status Code: %s", statusCodeConv))
	} else {
		OutputWrapper(ipAddrs, params, streams)
	}
	if args.AnalyzeHeader {
		AnalyzeHeaderWrapper(&consoleOutput, ipAddrsOut, client, params)
	} else {
		if ipAddrsOut != "" {
			consoleOutput.WriteString(fmt.Sprintf("\n\t╚═[ %s\n", ipAddrsOut))
		}
	}
	if args.PortScan != "" {
		PortScanWrapper(&consoleOutput, params, args)
	}
	utils.GJsonResult.Subdomains = append(utils.GJsonResult.Subdomains, utils.GSubdomBase)
	// Display the final result block
	fmt.Fprintln(utils.GStdout, consoleOutput.String())
	utils.GDisplayCount++
}
