# mantine-form-zod-resolver

[zod](https://www.npmjs.com/package/zod) resolver for [@mantine/form](https://mantine.dev/form/use-form/).

## Installation

With yarn:

```sh
yarn add zod mantine-form-zod-resolver
```

With npm:

```sh
npm install zod mantine-form-zod-resolver
```

## Zod versions

`mantine-form-zod-resolver` supports zod v3 and v4:

```tsx
// For zod v3
import { zodResolver } from 'mantine-form-zod-resolver';

// For zod v4
import { zod4Resolver } from 'mantine-form-zod-resolver';
```

## Basic fields validation

```tsx
import { z } from 'zod';
import { useForm } from '@mantine/form';
import { zodResolver } from 'mantine-form-zod-resolver';

const schema = z.object({
  name: z.string().min(2, { message: 'Name should have at least 2 letters' }),
  email: z.string().email({ message: 'Invalid email' }),
  age: z.number().min(18, { message: 'You must be at least 18 to create an account' }),
});

const form = useForm({
  initialValues: {
    name: '',
    email: '',
    age: 16,
  },
  validate: zodResolver(schema),
});

form.validate();
form.errors;
// -> {
//  name: 'Name should have at least 2 letters',
//  email: 'Invalid email',
//  age: 'You must be at least 18 to create an account'
// }
```

## Nested fields validation

```tsx
import { z } from 'zod';
import { useForm } from '@mantine/form';
import { zodResolver } from 'mantine-form-zod-resolver';

const nestedSchema = z.object({
  nested: z.object({
    field: z.string().min(2, { message: 'Field should have at least 2 letters' }),
  }),
});

const form = useForm({
  initialValues: {
    nested: {
      field: '',
    },
  },
  validate: zodResolver(nestedSchema),
});

form.validate();
form.errors;
// -> {
//  'nested.field': 'Field should have at least 2 letters',
// }
```

## List fields validation

```tsx
import { z } from 'zod';
import { useForm } from '@mantine/form';
import { zodResolver } from 'mantine-form-zod-resolver';

const listSchema = z.object({
  list: z.array(
    z.object({
      name: z.string().min(2, { message: 'Name should have at least 2 letters' }),
    })
  ),
});

const form = useForm({
  initialValues: {
    list: [{ name: '' }],
  },
  validate: zodResolver(listSchema),
});

form.validate();
form.errors;
// -> {
//  'list.0.name': 'Name should have at least 2 letters',
// }
```

## API:

### ZodResolverOptions

`zodResolver` takes as an optional second parameter some `zodResolverOptions`.

| Name            | Type                           | Description                                                                                                                                                                                                                        |
| --------------- | ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `errorPriority` | `first` \| `last` \| undefined | In case a field can display multiple error message, set `errorPriority` to `first` to display the message of the first failing check, or set `errorPriority` to `last` to display the message of the last failing check (default). |

## License

MIT
