package types

type Hostnames []Hostname
type Hostname string

type Stackname string

type AbsoluteDirs []AbsoluteDir

type RelativePath string
type AbsoluteDir AbsoluteEntry
type AbsoluteFilepath AbsoluteEntry
type AbsoluteEntry string

type CacheKey string

type ServiceType string //@TODO Merge with StackRole?
type ProgramName string
type Orgname string
type Version string

type Stacknames []Stackname

type Authorities []AuthorityDomain
type AuthorityDomain string

type Revision string

type StackRole string

type Basepath AbsoluteDir

type UrlTemplates []UrlTemplate
type UrlTemplate string

type RouteName string

type ResourceType string

type ResponseType string
