class Token:
  '''
  Information about an authentication token

  Properties:
  * name {string} token name
  * expire_date {Date}
  * secret_str {string} secret in string form suitable to add to HTTP request header
  '''
  def __init__(self, name, expire_date, secret_str):
    self.name = name
    self.expire_date = expire_date
    self.secret_str = secret_str


class CommonsConfig:
  '''
  Immutable configuration related to a specific data commons

  Properties:
  * name {string} arbitrary commons name to distinguish between commons
  * end_point {string} https:// url string of commons' API base -
       ex: https://dev.bionimbus.org/api/
  * refresh_token {Token} that can be used to acquire a session key
  * session_token {Token} last session key acquired if any
  '''
  def __init__(self, name, end_point, refresh_token, session_token):
    '''
    Constructor

    Arguments:
    * name {string}
    * end_point {string url}
    * refresh_token {Token}
    * session_token {Token}
    '''
    self.name = name
    self.end_point = end_point
    self.refresh_token = refresh_token
    self.session_token = session_token

  @staticmethod
  def load_from_file(name, path=None):
    '''
    Static factory that loads a CommonsConfig from a file

    Arguments:
    * name {string} name of data commons to load config for
    * path {string} defaults to ~/.cdis/name.json
    ''' 
    return CommonsConfig("dev","https://dev.bionimbus.org/api", None, None)
  
  def save_to_file(self, path):
    '''
    Serialize this config to a file at the given path

    Arguments:
    * path {string} defaults to ~/.cdis/(self.name).json
    '''
    raise NotImplementedError


class MetaNode:
  '''
  Generic metadata node

  Properties:
  * dict_type {string}
  * id {string uuid}
  * submitter_id {string}
  * project_id {string uuid}
  * project_submitter_id {string uuid}
  * props {dictionary} holds all the properties that belong to this dict_type
  '''
  def __init__(self, dict_type):
    self.dict_type = dict_type
  
  def set_kv(self, key, value):
    '''
    Helper throws ValidationError if invalid key/value specified for
    this node's dict_type
    '''
    raise NotImplementedError

  def get_v(self, key):
    raise NotImplementedError
  
  def validate(self):
    '''
    Return true if this node is valid to submit to the metadata store
    '''
    raise NotImplementedError


class IndexdEntry:
  '''
  Entry in object store index

  Properties:
  * id {string uuid}
  * path_list {list<uuid-string>}
  * attributes {dictionary<string,string>}
  '''
  def __init__(self, id, path_list, attributes):
    self.id = id


class ResultPage:
  '''
  Supports paging of fetch_meta_by_type results

  Properties:
  * data {Array<MetaNode>}
  * next_page may be None
  '''
  def __init__(self):
    '''
    '''
    raise NotImplementedError

class MetaFilter:
  '''
  Filter for CommonsHelper.list_meta_by_type which restricts
  the result to descendants of nodes with properties
  that match all of the key-valuelist filters

  Properties:
  * dict_type
  * property_to_values property name to value set
  '''
  def __init__(self, dict_type, property_to_values):
    self.dict_type = dict_type


class CommonsHelper:
  '''
  Helper class for some common data commons operations

  Properties:
  config {CommonsConfig} the currently active config
  dictionary {Dictionary} 
  '''
  def __init__(self, config):
    '''
    Client should create a CommonsHelper instance via the openSession factory method

    Arguments:
    * config {CommonsConfig}
    '''
    self.config = config
  
  @staticmethod
  def open_session(config, opts=None):
    '''
    Establish an API session with the commons backend.

    If config.session_token has not expired, then just use that, otherwise
    use config.refresh_token to open a new session.
    If ! helper.config.session_token is config.session_token, 
    then a new session has been created, and the client may
    want to helper.config.saveToFile()
    
    Also loads the dictionary from cache or the backend if opts['forceUpdate'] is specified or
    the dictionary is not already locally cached for config.end_point. 

    Returns:
    helper  
    '''
    raise NotImplementedError

  def fetch_meta_by_id(self, uuid_str):
    '''
    Given a uuid return a MetaNode or None if does not exist
    '''
    raise NotImplementedError
  
  def fetch_meta_by_name(self, project_submitter_id, submitter_id):
    raise NotImplementedError

  def fetch_meta_by_type(self, type_list, type_to_filter):
    '''
    Fetch the metadata nodes of the given types (in type_list) that match the given filters
    
    Ex:
        page1 = fetch_meta_by_type(['aliquot', 'variantFile'], { 
          'project' : { 'submitter_id': ['whatever'] },
          'case' : { 'species': ['human'] },
          'demographic': { 'race': [ 'black', 'african american'], 'year_of_birth': [ '< 1980' ] }
        })

    1. For each type in type_to_filter, find every node of that type that matches the given filters
        where dict_type is implicitly in type_to_filter if not explicitly there
    2. For every node in 1 compute its path to the metadata root
    3. Compute the intersection of all the paths in 2

    The result includes every node of dict_type with a path to a root project node
    that intersects the path-to-root of every node 

    Arguments:
    * dict_type type of nodes to fetch
    * type_to_filter mapping of type to MetaFilter
    Return:
    The first ResultPage of meta data nodes of the given dict_type that match the given filters
    '''
    raise NotImplementedError

  def fetch_next_page(self, last_page):
    raise NotImplementedError
  
  def new_meta_node(self, dict_type):
    '''
    Create a MetaNode initialized with default values for the given dictionary type

    Returns:
    new MetaNode ready for the client to initialize with property values
    '''
    raise NotImplementedError
  

  def fetch_indexd_entry(self, uuid):
    '''
    Fetch the object store index entry associated with the file with the given uuid -
    the uuid may be discovered by querying the metadata db via fetchNodeX or findNodesX
    '''
    raise NotImplementedError
  
  def upload_file(self, uuid, path):
    '''
    Fetch the meta_node with uuid, and
    upload the file at path and associated with the given meta_node
    to the commons object store, and update the object store index (indexd) accordingly.

    Arguments:
    * path {string}
    * meta_node {MetaNode}

    Return:
    A Promise<IndexdEntry> that resolves after the file has been successfully posted 
    to the object store and indexd
    '''
    raise NotImplementedError

  def download_file(self, uuid, path):
    '''
    Download the file with uuid to path (path should not already exist)

    Return:
    A Promise<IndexdEntry> that resolves after the file has been downloaded to path
    '''
    raise NotImplementedError
  
  def save_meta_nodes(self, node_list):
    '''
    Arguments:
    * node_list {List<MetaNode>}

    Return: 
    A Promise/Future that resolves or fails when the save has completed.
    On success the promise resolves to the list of nodes as saved on the backend.
    On failre the promise rejects to the list of nodes with errors List<Tuple<MetaNode,ErrorString>> 
    '''
    raise NotImplementedError
  