# Gearbox


## Target User Types
- Sitebuilders
	- Install and configure sites
- Software Engineers
	- Implement sites with custom programming
- GearBuilders 
	- Create _"Gears"_ for Gearbox


## APIs
- Host API - Used for host things, before Box is running
- Box API - Used for Box things

## Gearbox Components
- "Single" Executable
- GearboxOS ISO _(Based on Alpine)_		

## GearboxOS
- Restful API - Interact with Admin
- Built-in Go-based SSH
- Broadcast `.local` hostnames
	- NO EDITING HOSTS FILE!!!!


## "Single" Executable
- Compiled #GoLang code
	- Same source code for all variants
	- One executable for a host computer 
	    - On a Windows `gearbox.exe`
	    - On a Mac, Linux `gearbox`
	- One executable inside the Gearbox Box
	    - `gearbox`


## On First Execution
- Start Admin console from local host
- Check for VirtualBox
	- If installed ensure correct version and configuration
	- If not installed, install and configure if not installed
- Download ISO for GearboxOS
- Create Box using VirtualBox and the ISO 
- Configure Box correctly
	
## On Every Execution
- Start Admin console from local host
- Re-verify VirtualBox &  Box
- Loads User Config file 
- Scans	Project(s) directories
	- Looks for:
		- Subdirs with `project.json` file => Projects
			- When new project is found
			- Set to disabled
		- Other subdirs => Project candidates

## Projects
- Enabled
	- User Config marks as `"active"`
	- Causes Gearbox to run required stack of service containers
		- If project requires Nginx, an Nginx will be run
- Disabled
- Candidates
	- Has no project.json
	- Turn into a Disabled or Enabled project 
		- By adding `project.json`
	- Recognize VVV or Local or DDEV or Lando box and allow import

## Gearbox configuration
- User config directory (Gearbox Concept)		
- macOS: 

## Collection of projects
- For each Project 
	- Project Config: `gearbox.json`
	
## Go Tools
- Regenerate assets.go 
	- https://github.com/jteeuwen/go-bindata
	- `go-bindata -dev -o assets.go -pkg gearbox admin/...`