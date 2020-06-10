
func k8sStartFuncDemo() *cobra.Command{
    s := options.NewServerRunOptions()
    
    cmd := &cobra.Command{
        ...
        completedOptions, err := Complete(s)
        return Run(completedOptions, stopCh)
    }

    fs := cmd.Flags()
    // handler flags

    cmd.SetUsageFunc()
    cmd.SetHelpFunc()

    return cmd
}
