import time
from datetime import datetime

from cdis.gen3 import sdk

test_project = 'test_project'

def get_helper():
  config = sdk.CommonsConfig.load_from_file('test')
  helper = sdk.CommonsHelper.open_session(config)
  return helper

def test_get_folder():
  '''
  Tests go into dated folders
  '''
  helper = get_helper()
  now = datetime.now()
  submitter_id = 'test_' + str(time.time())
  test_folder = helper.fetch_meta_by_name(test_project, submitter_id)
  if test_folder is None:
    test_folder = helper.new_meta_node('test_folder')
    test_folder.project_code = test_project;
    test_folder.submitter_id = submitter_id
    result = helper.save_meta_nodes([ test_folder ]).wait()
    assert result.is_ok
    test_folder = result.node_list[0]
  assert test_folder.get_v('submitter_id') is submitter_id
  assert test_folder.dict_type is 'test_folder'
  return test_folder
  
def test_list_projects():
  helper = get_helper()
  page1 = helper.fetch_meta_by_type(['project'])
  assert page1.next_page is None
  assert page1.data.length > 0


def test_submit_test_file():
  helper = get_helper()
  test_folder = test_get_folder()
  test_file = helper.new_meta_node('test_file')
  submitter_id = 'test_' + str(time.time())
  test_file.set_kv('submitter_id', submitter_id)
  test_file.set_kv('project_id', test_project)
  test_file.set_kv('test_folder.submitter_id', test_folder.get_v('submitter_id'))
  result = helper.save_meta_nodes([ test_file ]).wait(10, lambda: println('waiting ...'))
  assert result.is_ok
  test_file = result.node_list[0]
  result = helper.upload_file(test_file.get_v('uuid'), './test_file.txt').wait(10, lambda: println('uploading ...'))
  assert result.is_ok
  result = helper.download_file(test_file.get_v('uuid'), '/tmp/' + submitter_id).wait(10, lambda: println('downloading ...'))
  assert result.is_ok
  result = helper.fetch_meta_by_type([ 'test_file' ], { 'test_file': [('submitter_id',submitter_id)] }).wait(
    10, lambda: println('searching ...')
  )
  assert result.is_ok
  assert result.data.length is 1
  assert result.data[0].dict_type is 'test_file'
  assert result.data[0].get_v('submitter_id') is submitter_id

def test_fetch_by_type():
  helper = get_helper()
  result = helper.fetch_meta_by_type([ 'test_file' ], 
      { 
        'test_file': [ ('submitter_id',submitter_id) ],
        'test_folder': [ ('create_date', ':after 01/01/2017') ],
        'demographic': [ ('country', 'United States'), ('age', '> 40'), ('age', '< 50'), ('ethnicity', 'white') ] 
      }
    ).wait(
      10, lambda: println('searching ...')
    )
  assert result.is_ok
  assert result.data.length > 0

  