import { SnapshotTarget } from "./store";

export const addCollectJob = async (
  endpoint: string,
  target: SnapshotTarget
): Promise<void> => {
  const resp = await fetch(`/api/${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(target),
  });

  if (!resp.ok) {
    return alert(
      `http error: status=${resp.status}, message=${await resp.text()}`
    );
  }
};
