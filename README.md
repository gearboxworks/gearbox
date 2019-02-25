# Gearbox


## Tech Stack
- BEM      - http://getbem.com/naming/
- CSS-Grid - 2d 
- Flexbox  - 1d 

## First Admin Milestone
- Create a mockup for approval
- Create HostConnector for Windows
- Ask to change project dir if not projects found
- Display enabled and disable projects
- Allow changing project dir
- Allow adding additional project dir
- Allow displaying candidate projects
- Allow enabling a project


## Next Gearbox Milestone
- `gearbox vm start`
- `gearbox vm stop`
- `gearbox ssh`


## Target User Types
- Sitebuilder
	- Install sites
- Engineers
	- Engineer sites
- GearBuilders 
	- Create containers


## APIs
- Host API - Used for host things, before VM is running
- VM API - Used for VM things

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
	- One executable inside the Gearbox VM
	    - `gearbox`


## On First Execution
- Start Admin console from local host
- Check for VirtualBox
	- If installed ensure correct version and configuration
	- If not installed, install and configure if not installed
- Download ISO for GearboxOS
- Create Gearbox VM using VirtualBox and the ISO 
- Configure Gearbox VM correctly
	
## On Every Execution
- Start Admin console from local host
- Reverify VirtualBox & Gearbox VM
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
	- Project Config: `project.json`
	
## Go Tools
- Regenerate assets.go 
	- https://github.com/jteeuwen/go-bindata
	- `go-bindata -dev -o assets.go -pkg gearbox admin/...`