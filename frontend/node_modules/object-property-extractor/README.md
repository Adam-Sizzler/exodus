# Object Property Extractor

Access deep object properties using a string (e.g. "user.country.name")

A lightweight (no dependencies) tool to extract deeply nested values from JS Objects (or Arrays), with optional Fallback.

Similar to [Lodash' `get`](https://lodash.com/docs/4.17.15#get) function, but with some additional functionality.

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

In Javascript, you call inner object properties via dot notation:
```js
data.user.name.last // Fett
```

If you want to access a property dynamically, you can do this:
```js
const key = "user" 
return data[key]
```

However, you *can't* do this:
```js
const key = "user.name"
return data[key]
```

This tool allows access to deep properties from a single "property path" string.

## Installation

```js
yarn add object-property-extractor
// OR
npm install object-property-extractor
```

## Usage

`extract( dataObject, propertyString, [fallback] )`

```js
import extract from "object-property-extractor"

// Using the data object above
extract(data, "user.name.first") // Jango

// With fallback when path not found
extract(data, "user.age", "Unknown") // Unknown

// Arrays can be accessed by index, as per normal indexing syntax
extract(data, "user.children[1]") // Boba
```

### Array handling

In addition to accessing array by index (above), if an array consists of objects, then it's possible to extract a single property from *each* object in the returned array.

For example:
```js
extract(data, "user.weapons.name")
// ["Blaster", "Seismic charge"]
```

Note that this is essentially a shorthand for:  
`extract(data, "user.weapons").map((weapon) => weapon.name)`

### Error handling

If a requested property can't be accessed (e.g. incorrect path), the function will *throw* an error, unless a fallback is provided. So unless you are catching and handling these errors at a higher level, it is recommended to always provided a fallback (`null` is an acceptable fallback).

## Testing

A [jest](https://jestjs.io/) test suite is included in the repo. To run:
`yarn test`

See `/test/test.ts` for the test cases.

## Bug report / Feature requests

Please make an issue in the Github repo: https://github.com/CarlosNZ/object-property-extractor
