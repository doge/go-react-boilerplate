import { Toaster } from "sonner";
import { useTheme } from "../../providers/theme-provider";

export function ThemedToaster() {
  const { theme } = useTheme();
  return <Toaster richColors position="bottom-center" theme={theme} />;
}
