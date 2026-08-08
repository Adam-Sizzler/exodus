const METHOD_COLORS = new Map<string, string>([
    ['GET', 'teal'],
    ['POST', 'blue'],
    ['PUT', 'orange'],
    ['PATCH', 'yellow'],
    ['DELETE', 'red'],
    ['OPTIONS', 'grape'],
    ['HEAD', 'indigo']
])

export const getMethodColor = (method: string): string => METHOD_COLORS.get(method) ?? 'gray'
