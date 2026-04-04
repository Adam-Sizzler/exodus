# Object Property Assigner

*Companion package to [object-property-extractor](https://www.npmjs.com/package/object-property-extractor)*

A lightweight (no dependencies) tool to assign deeply nested values in JS Objects (or Arrays).

## Why?

Consider the object
```js
const data = {
  user: {
    name: { first: 'Jango', last: 'Fett' },
    children: ['Boba', 'Clone 1', 'Clone 2', ...etc],
    weapons: [
      { name: 'Blaster', description: 'For shooting stuff' },
      { name: 'Seismic charge', description: '...BWAAAAAANG' },
    ],
  },
  ...otherProperties,
}
```

In Javascript, you assign inner object properties via dot notation:
```js
data.user.name.last = "Hutt" 
```

If you want to assign a property dynamically, you can do this:
```js
const key = "user" 
data[key] = { name: "Boba Fett" } // data.user = { name: "Boba Fett" }
```

However, you *can't* do this:
```js
const key = "user.name"
data[key] = "Boba Fett"
```

This tool allows assignation of deep properties using a single "property path" string.

## Installation

```js
yarn add object-property-assigner
// OR
npm install object-property-assigner
```

## Usage

*Note: function returns a shallow copy of the input, the original object is unmodified*

`assign( dataObject, propertyString, newValue, { options } )`
(See below for [`options` details](#options))

```js
import assign from "object-property-assigner"

// Using the data object above
data = assign(data, "user.name.first", "Boba") // data.user.name.first = "Boba"

data = assign(data, "user.weapons[1].description", "Pew Pew") // data.user.weapons[1].description = "Pew Pew"
```

### Array handling

In addition to accessing array by index (above), if an array consists of objects, then it's possible to assign a single property for *all* object in the array.

For example:
```js
data = assign(data, "user.weapons.name", "Laser Gun")
// sets *all* user.weapons.name to "Laser Gun"
```


### Options

The (optional) `options` object can contain any or all of the following parameters:

- `remove` -- if `true`, the property will be *deleted* rather than assigned (in which case the `newValue` parameter is ignored)  
  ```js
  assign(data, "user.name.first", null, {remove:true}) // delete user.name.first
  ```
- `createNew` -- (default: `true`). If a property doesn't exist, it will be created, so set this to `false` if this behaviour is not desired  
  ```js
  assign(data, "user.kind", "Mandalorian", {createNew: true}) // data.user.kind = "Mandalorian"
  ```  
  Note: for arrays, if an index is specified higher than the current array length, a new item will be created as the *next* array item, regardless of how much higher the index is. e.g.  
  ```js
  assign( {myArray: [1, 2, 3]}, "myArray[10]", "New Value")
  // --> {myArray: [1, 2, 3, "New Value"]}
  ```
- `insert` (`boolean`) -- (only for arrays). If `true`, then the new value will be *inserted* into the array at the designated position and all other elements shuffled along. e.g.  
  ```js
  assign( {myArray: [1, 2, 3]}, "myArray[1]", "New Value",
    {insert: true} 
  )
  // --> {myArray: [1, "New Value", 2, 3]}
  ```
- `insertBefore` (`string | number`) -- (only for objects). If this property is set, then the new value will be inserted *before* the key specified in `insertBefore`. e.g.  
  ```js
  assign( {one: 1, two: 2, three: 3}, "myNewValue", 69,
    {insertBefore: "three"}
  )
  // --> {one: 1, two: 2, myNewValue: 69, three: 3}
  ```
  If value is a number, it will be considered an *index* and treat the current properties as an array.

  **Please note**: Javascript does *not* guarantee object property order, so this option may yield unpredictable results, especially if using non-standard object keys (e.g. numbers). See [here](https://dev.to/frehner/the-order-of-js-object-keys-458d) for an explanation of how key ordering is handled.
- `insertAfter` (`string`) -- (only for objects). Same as `insertBefore`, but inserts *after* a particular object key. Same caveats apply.
- `noError` -- (default: `false`). If a property doesn't exist *and* `createNew == false`, then an error will be thrown. If you'd rather it just silently ignored the missing property, then set this parameter to `true`. Note that this is only for errors due to invalid property strings -- other errors might still be thrown.

## Testing

A [jest](https://jestjs.io/) test suite is included in the repo. To run:
`yarn test`

See `/test/test.ts` for the test cases.

## Bug report / Feature requests

Please make an issue in the Github repo: https://github.com/CarlosNZ/object-property-assigner
