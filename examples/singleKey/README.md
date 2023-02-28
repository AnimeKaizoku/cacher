# Single Key System
This type of system contains a single element as a key.

It has been divided into 3 types:
1. built-in type key vs built-in type value (located at [examples/singleKey/builtin-builtin](examples/singleKey/builtin-builtin))
2. built-in type key vs custom type value (located at [examples/singleKey/builtin-custom](examples/singleKey/builtin-custom))
3. custom type key vs custom type value (located at [examples/singleKey/custom-custom](examples/singleKey/custom-custom))

For the sake of simplicity we will be using practical life contexts in these examples as described below:

  - Type 1 will be a int64-string mapping, where int64 key would be chat id of our
    hypothetical chatting program, and string value would be the title of that chat.

  - Type 2 will be a int-struct mapping, where int64 key would be user id of our
    hypothetical user for some storage service, and struct value would contain user's
	   name (type string), dob (type struct{date int, month int, year int}).
  
  - Type 3 will be a (custom int)-([]struct)  mapping, where custom int will be 
    rank level of a person in an agency from (1, 2, 3) and value array will contain
    struct consisting the user's name and dob.

**Note**: You should try out experiment mentioned in introductory docstrings of Type 2
single key system.