#ifndef TRAY_H
#define TRAY_H

struct tray_menu;

struct tray {
    struct {
        char* data; // Will be real data on darwin and windows, but path on linux
        int length;
    } icon;
    struct tray_menu *menu;
};

struct tray_menu {
    char *text;
    int disabled;
    int checked;

    void (*cb)(struct tray_menu *);
    void *context;

    struct tray_menu *submenu;
};


int tray_init(struct tray *tray);
void tray_update(struct tray *tray);
int tray_loop(int blocking);
void tray_exit();

#endif /* TRAY_H */
