type DateFormatStyle = "short-date" | "date-time";

const DATE_FORMAT_OPTIONS: Record<DateFormatStyle, Intl.DateTimeFormatOptions> =
  {
    "short-date": {
      day: "2-digit",
      month: "short",
      year: "numeric",
    },
    "date-time": {
      day: "2-digit",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    },
  };

export function formatDisplayDate(
  value?: string,
  style: DateFormatStyle = "date-time",
  fallback = "Tanggal belum tersedia",
): string {
  if (!value) {
    return fallback;
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return fallback;
  }

  return new Intl.DateTimeFormat("id-ID", DATE_FORMAT_OPTIONS[style]).format(
    date,
  );
}
