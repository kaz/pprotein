import { SnapshotTarget } from "./store";

export const generateId = new Date().getTime().toString(36);

export const addCollectJob = async (
  endpoint: string,
  target: SnapshotTarget
) => {
  const resp = await fetch(`/api/${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(target),
  });

  if (!resp.ok) {
    alert(await resp.text());
  }
};
