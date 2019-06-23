package types

type Name = string

type Hostnames []Hostname
type Hostname = string

type Stackname = string

type Dirs []Dir

type Path = string
type Dir = FileSystemEntry
type Filepath = FileSystemEntry
type FileSystemEntry = string

type CacheKey string

type ServiceType = string //@TODO Merge with StackRole?
type ProgramName = string
type Orgname = string
type Version = string

type Stacknames []Stackname

type AuthorityDomains []AuthorityDomain
type AuthorityDomain = string

type Revision = string

type StackRole = string

type Basepath = Dir

type UrlTemplates []UrlTemplate
type UrlTemplate = string

type RouteName = string

type ResourceType = string

type ResponseType = string
