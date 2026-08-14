import { useState } from "react";
import { avatarColor, nameInitials } from "../utils/avatar";

export type AvatarSize = "xs" | "sm9" | "sm" | "lg";

const SIZE_CLASSES: Record<AvatarSize, string> = {
  xs: "h-8 w-8 text-xs",
  sm9: "h-9 w-9 text-xs",
  sm: "h-10 w-10 text-sm",
  lg: "h-28 w-28 text-4xl",
};

type AvatarProps = {
  name: string;
  avatarUrl?: string | null;
  size?: AvatarSize;
  title?: string;
  className?: string;
};

export default function Avatar({ name, avatarUrl, size = "sm", title, className = "" }: AvatarProps) {
  const [broken, setBroken] = useState(false);
  const sizeClasses = SIZE_CLASSES[size];

  if (avatarUrl && !broken) {
    return (
      <img
        src={avatarUrl}
        alt={name}
        title={title}
        className={`${sizeClasses} rounded-full object-cover shrink-0 ${className}`}
        onError={() => setBroken(true)}
      />
    );
  }

  return (
    <div
      title={title}
      className={`${sizeClasses} rounded-full flex items-center justify-center text-white font-semibold shrink-0 ${avatarColor(name)} ${className}`}
    >
      {nameInitials(name)}
    </div>
  );
}
