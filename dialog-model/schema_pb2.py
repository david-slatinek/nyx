# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: schema.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0cschema.proto\x12\x06schema\"\x1d\n\rDialogRequest\x12\x0c\n\x04text\x18\x01 \x01(\t\" \n\x0e\x44ialogResponse\x12\x0e\n\x06\x61nswer\x18\x01 \x01(\t2H\n\rDialogService\x12\x37\n\x06\x44ialog\x12\x15.schema.DialogRequest\x1a\x16.schema.DialogResponseB\x04Z\x02./b\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'schema_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\002./'
  _DIALOGREQUEST._serialized_start=24
  _DIALOGREQUEST._serialized_end=53
  _DIALOGRESPONSE._serialized_start=55
  _DIALOGRESPONSE._serialized_end=87
  _DIALOGSERVICE._serialized_start=89
  _DIALOGSERVICE._serialized_end=161
# @@protoc_insertion_point(module_scope)
