package typescript

import "ella.to/transform"

const helperFunc = `
// Helper Functions

interface Subscription<EventMap> {
  addEventListener<K extends keyof EventMap>(
    type: string | number | symbol,
    listener: (event: EventMap[K]) => void
  ): void;
  close(): void;
}

async function callServiceStreamMethod<Req, EventMap>(
  host: string,
  path: string,
  method: "GET",
  body?: Req,
  headers?: Record<string, string>
): Promise<Subscription<EventMap>> {
  if (method !== "GET") {
    throw new Error("only GET method is supported for streaming");
  }

  const url = createURL(host, path, prepareForQs(body));
  const sse = new EventSource(url, { withCredentials: true, headers } as any);

  return new Promise((resolve, reject) => {
    sse.onerror = (e) => {
      reject(e);
    };

    sse.onopen = () => {
      resolve({
        addEventListener<K extends keyof EventMap>(
          event: K,
          listener: (event: EventMap[K]) => void
        ) {
          sse.addEventListener(event as any, listener as any);
        },
        close() {
          sse.close();
        },
      });
    };
  });
}

async function callServiceMethod<Req, Resp>(
  host: string,
  path: string,
  method: "GET" | "POST",
  body?: Req,
  headers?: Record<string, string>
) {
  const url =
    method === "GET"
      ? createURL(host, path, prepareForQs(body))
      : createURL(host, path);

  if (method === "GET") {
    body = undefined;
  }

  const resp = await fetch(url, {
    method: method,
    body: body ? JSON.stringify(body) : undefined,
    headers: headers,
    credentials: "include",
  });

  const value = await resp.text();

  if (resp.status > 300) {
    try {
      const msg = JSON.parse(value);
      throw new Error(msg.error);
    } catch (e) {
      throw new Error(value);
    }
  }

  return JSON.parse(value) as Resp;
}

function prepareForQs(obj?: any): Record<string, string> | undefined {
  if (!obj) {
    return undefined;
  }

  const record: Record<string, string> = {};
  for (const key in obj) {
    const value = obj[key];
    if (
      typeof value !== "string" &&
      typeof value !== "number" &&
      typeof value !== "boolean"
    ) {
      throw new Error("Invalid value type for key: " + key);
    }
    record[key] = encodeURIComponent(obj[key] + "");
  }
  return record;
}

function createURL(
  host: string,
  path: string,
  qs?: Record<string, string>
): string {
  const url = new URL(host);
  url.pathname = (url.pathname + path).replace(/\/\//g, "/");
  if (qs) {
    for (const key in qs) {
      url.searchParams.append(key, qs[key]);
    }
  }
  return url.href;
}
`

func HelperFunc() transform.Func {
	return func(out transform.Writer) error {
		out.Str(helperFunc)
		return nil
	}
}
